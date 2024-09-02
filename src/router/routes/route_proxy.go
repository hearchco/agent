package routes

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/hearchco/agent/src/search/useragent"
	"github.com/hearchco/agent/src/utils/anonymize"
	"github.com/hearchco/agent/src/utils/moreurls"
)

func routeProxy(w http.ResponseWriter, r *http.Request, salt string, timeout time.Duration) error {
	// Parse the form.
	err := r.ParseForm()
	if err != nil {
		// Server error.
		werr := writeResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse form: %v", err))
		if werr != nil {
			return fmt.Errorf("%w: %w", werr, err)
		}
		return err
	}

	// Get the parameters.
	params := r.Form
	urlParam := getParamOrDefault(params, "url")
	hashParam := getParamOrDefault(params, "hash")
	faviconParam := getParamOrDefault(params, "favicon", strconv.FormatBool(false))

	// Check the required parameters.
	if urlParam == "" || hashParam == "" {
		// User error.
		return writeResponse(w, http.StatusBadRequest, "url and hash are required")
	}

	// Check if only favicon is requested.
	favicon, err := strconv.ParseBool(faviconParam)
	if err != nil {
		// User error.
		return writeResponse(w, http.StatusBadRequest, "favicon must be a boolean")
	}

	// Get url to verify and to proxy.
	urlToVerify, urlToProxy, err := getUrlToVerifyAndToProxy(urlParam, favicon)
	if err != nil {
		// User error.
		log.Debug().
			Err(err).
			Str("url", urlParam).
			Str("url_to_verify", urlToVerify).
			Str("url_to_proxy", urlToProxy).
			Str("hash", hashParam).
			Str("favicon", faviconParam).
			Msg("Failed to get URL to verify and to proxy")
		return writeResponse(w, http.StatusBadRequest, fmt.Sprintf("failed to get URL to verify and to proxy: %v", err))
	}

	// Check if hash is valid.
	if !anonymize.VerifyHMACBase64(hashParam, urlToVerify, salt) {
		// User error.
		log.Debug().
			Str("url", urlParam).
			Str("url_to_verify", urlToVerify).
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
			Str("url", urlParam).
			Str("url_to_proxy", urlToProxy).
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
	log.Debug().
		Str("url", target.String()).
		Msg("Proxying request")

	// Use the new request.
	resp, err := requestResponse(&nr, timeout)
	if err != nil {
		// Server error.
		log.Debug().
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

func getUrlToVerifyAndToProxy(urlParam string, favicon bool) (string, string, error) {
	urlToVerify := urlParam
	urlToProxy := urlParam

	if favicon {
		// Get the URI to verify.
		urlUri, err := moreurls.GetURIToVerify(urlParam)
		if err != nil {
			return "", "", fmt.Errorf("failed to extract URI from URL: %w", err)
		}

		// Get the favicon URL.
		faviconUrl, err := getFaviconURL(urlParam)
		if err != nil {
			return "", "", fmt.Errorf("failed to extract favicon URL: %w", err)
		}

		// Set the URLs.
		urlToVerify = urlUri
		urlToProxy = faviconUrl
	}

	return urlToVerify, urlToProxy, nil
}

// Appends the path to favicon to the URI of the URL.
func getFaviconURL(urll string) (string, error) {
	// TODO: Impl getting the favicon path from the html head.
	const faviconPath = "/favicon.ico"
	uri, err := moreurls.GetURI(urll)
	if err != nil {
		return "", err
	} else {
		return uri + faviconPath, nil
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
