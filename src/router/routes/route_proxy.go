package routes

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strconv"

	"github.com/rs/zerolog/log"

	"github.com/hearchco/agent/src/config"
	"github.com/hearchco/agent/src/search/useragent"
	"github.com/hearchco/agent/src/utils/anonymize"
)

func routeProxy(w http.ResponseWriter, r *http.Request, salt string, timeouts config.ImageProxyTimeouts) error {
	err := r.ParseForm()
	if err != nil {
		// Server error.
		werr := writeResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse form: %v", err))
		if werr != nil {
			return fmt.Errorf("%w: %w", werr, err)
		}
		return err
	}

	params := r.Form

	urlParam := getParamOrDefault(params, "url")
	hashParam := getParamOrDefault(params, "hash")
	faviconParam := getParamOrDefault(params, "favicon", strconv.FormatBool(false))

	if urlParam == "" || hashParam == "" {
		// User error.
		return writeResponse(w, http.StatusBadRequest, "url and hash are required")
	}

	urlToProxy := urlParam
	if val, err := strconv.ParseBool(faviconParam); err != nil {
		// User error.
		return writeResponse(w, http.StatusBadRequest, "favicon must be a boolean")
	} else if val {
		faviconUrl, err := getFaviconURL(urlParam)
		if err != nil {
			// User error.
			log.Debug().
				Err(err).
				Str("url", urlParam).
				Str("favicon", faviconUrl).
				Msg("Failed to create favicon URL")
			return writeResponse(w, http.StatusBadRequest, err.Error())
		} else {
			urlToProxy = faviconUrl
		}
	}

	// Check if hash is valid.
	if !anonymize.VerifyHash(hashParam, urlToProxy, salt) {
		// User error.
		log.Debug().
			Str("url", urlParam).
			Str("url_to_proxy", urlToProxy).
			Str("hash", hashParam).
			Str("favicon", faviconParam).
			Msg("Invalid hash")
		return writeResponse(w, http.StatusUnauthorized, "invalid hash")
	}

	// Parse the url.
	target, err := url.Parse(urlToProxy)
	if err != nil {
		// User error.
		log.Debug().
			Str("url", urlToProxy).
			Msg("Invalid url")
		return writeResponse(w, http.StatusBadRequest, "invalid url")
	}

	// Create a new request.
	nr := createAnonRequest(r, target)
	log.Trace().
		Caller().
		Str("request", fmt.Sprint(nr)).
		Msg("Created a new anon request")

	// Create reverse proxy with timeout.
	rp := createReverseProxy(timeouts)

	// Proxy the request.
	log.Debug().
		Str("url", target.String()).
		Msg("Proxying request")
	rp.ServeHTTP(w, &nr) // Use the new request.

	return nil
}

// Appends the path to favicon to the URI of the URL.
func getFaviconURL(urll string) (string, error) {
	// TODO: Impl getting the favicon path from the html head.
	const faviconPath = "/favicon.ico"
	uri, err := getURI(urll)
	if err != nil {
		return "", err
	} else {
		return uri + faviconPath, nil
	}
}

// Extracts the URI from the URL.
// https://www.example.com/some/path -> https://www.example.com
func getURI(urll string) (string, error) {
	const uriPattern = "^(http(s?))(://)([^/]+)"
	re := regexp.MustCompile(uriPattern)
	ss := re.FindString(urll)
	if ss == "" {
		return "", fmt.Errorf("failed to extract URI from URL")
	} else {
		return ss, nil
	}
}

func createAnonRequest(r *http.Request, target *url.URL) http.Request {
	// Get random UserAgent with corresponding Sec-Ch-Ua headers.
	ua := useragent.RandomUserAgentWithHeaders()

	// Create a new request.
	return http.Request{
		Method:     http.MethodGet,
		URL:        target,
		Host:       target.Host,
		RequestURI: target.RequestURI(),
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

func createReverseProxy(timeouts config.ImageProxyTimeouts) httputil.ReverseProxy {
	// Create reverse proxy with timeout.
	rp := httputil.ReverseProxy{Director: func(r *http.Request) {}}
	rp.Transport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   timeouts.Dial,
			KeepAlive: timeouts.KeepAlive,
		}).DialContext,
		TLSHandshakeTimeout: timeouts.TLSHandshake,
	}
	return rp
}
