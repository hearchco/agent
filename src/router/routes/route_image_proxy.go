package routes

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/hearchco/agent/src/search/useragent"
	"github.com/hearchco/agent/src/utils/anonymize"
)

func routeImageProxy(w http.ResponseWriter, r *http.Request, secret string, timeout time.Duration) error {
	// Parse the form.
	err := r.ParseForm()
	if err != nil {
		// Server error.
		log.Error().
			Caller().
			Err(err).
			Msg("Failed to parse form")
		return writeResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse form: %v", err))
	}

	// Get the parameters.
	params := r.Form
	urlParam := getParamOrDefault(params, "url")
	fqdnParam := getParamOrDefault(params, "fqdn")
	hashParam := getParamOrDefault(params, "hash")
	timestampParam := getParamOrDefault(params, "timestamp")

	// Check the required parameters.
	if (urlParam == "" && fqdnParam == "") || hashParam == "" || timestampParam == "" {
		// User error.
		return writeResponse(w, http.StatusBadRequest, "url/fqdn, hash and timestamp are required")
	}

	// Check if both url and fqdn are provided.
	if urlParam != "" && fqdnParam != "" {
		// User error.
		return writeResponse(w, http.StatusBadRequest, "only one of url or fqdn is allowed")
	}

	// Check wether to use url or fqdn.
	var favicon bool
	if fqdnParam != "" {
		favicon = true
	}

	// Get url to verify and to proxy.
	var verificator string
	if !favicon {
		verificator = urlParam
	} else {
		verificator = fqdnParam
	}

	// Check if hash is valid.
	if ok, err := anonymize.VerifyHMACBase64(hashParam, verificator, secret, timestampParam); !ok || err != nil {
		// User error.
		log.Debug().
			Err(err).
			Str("verificator", verificator).
			Str("hash", hashParam).
			Str("timestamp", timestampParam).
			Msg("Invalid hash")
		return writeResponse(w, http.StatusUnauthorized, "invalid hash")
	}

	// Parse the url.
	var target *url.URL
	if !favicon {
		urll, err := url.Parse(urlParam)
		if err != nil {
			// User error.
			log.Debug().
				Str("url", urlParam).
				Msg("Invalid url")
			return writeResponse(w, http.StatusBadRequest, "invalid url")
		}
		target = urll
	} else {
		faviconUrll := getFaviconURL(fqdnParam)
		urll, err := url.Parse(faviconUrll)
		if err != nil {
			// Server error.
			log.Error().
				Err(err).
				Str("fqdn", fqdnParam).
				Msg("Failed to get favicon url")
			return writeResponse(w, http.StatusInternalServerError, "failed to get favicon url")
		}
		target = urll
	}

	// Create a new request.
	nr := createAnonRequest(r, target)
	log.Trace().
		Caller().
		Str("request", fmt.Sprint(nr)).
		Msg("Created a new anon request")

	// Create reverse proxy with timeout.
	log.Debug().
		Str("url", target.String()).
		Msg("Proxying request")

	// Use the new request.
	resp, err := requestResponse(&nr, timeout)
	if err != nil {
		// Server error.
		log.Error().
			Err(err).
			Str("url", target.String()).
			Msg("Failed to proxy request")
		return writeResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to proxy request: %v", err))
	}

	log.Trace().
		Caller().
		Str("response", fmt.Sprint(resp)).
		Msg("Got a response")

	// Proxy the response.
	return writeResponseImageProxy(w, resp)
}

// Returns the url of the favicon for the provided FQDN.
func getFaviconURL(fqdn string) string {
	// TODO: Impl getting the favicon from favicon providers.
	const faviconPath = "favicon.ico"
	return fmt.Sprintf("https://%s/%s", fqdn, faviconPath)
}

func createAnonRequest(r *http.Request, target *url.URL) http.Request {
	// Get random UserAgent with corresponding Sec-Ch-Ua headers.
	ua := useragent.RandomUserAgentWithHeaders()

	// Create a new request.
	return http.Request{
		Method:     http.MethodGet,
		URL:        target,
		Host:       target.Host,
		Proto:      "HTTP/2",
		ProtoMajor: 2,
		ProtoMinor: 0,
		Header: map[string][]string{
			"Accept":             {"image/avif", "image/webp", "image/apng", "image/svg+xml", "image/*", "*/*;q=0.8"},
			"Accept-Encoding":    r.Header["Accept-Encoding"], // WARN: This is passed from the original request.
			"Accept-Language":    {"en-US,en;q=0.9"},
			"Sec-Ch-Ua":          {ua.SecCHUA},
			"Sec-Ch-Ua-Mobile":   {ua.SecCHUAMobile},
			"Sec-Ch-Ua-Platform": {ua.SecCHUAPlatform},
			"Sec-Fetch-Dest":     {"image"},
			"Sec-Fetch-Mode":     {"no-cors"},
			"Sec-Fetch-Site":     {"same-site"},
			"User-Agent":         {ua.UserAgent},
		},
	}
}

func requestResponse(r *http.Request, timeout time.Duration) (*http.Response, error) {
	// Create a reverse proxy.
	client := http.Client{Timeout: timeout}
	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	}

	// Check if the response is OK.
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status code is not OK: %v", resp.StatusCode)
	}

	// Check if the response is of image type.
	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "image/") {
		return nil, fmt.Errorf("response content type is not an image: %v", resp.Header.Get("Content-Type"))
	}

	return resp, nil
}
