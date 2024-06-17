package routes

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"

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

	if urlParam == "" || hashParam == "" {
		// User error.
		return writeResponse(w, http.StatusBadRequest, "url and hash are required")
	}

	// Check if hash is valid.
	if !anonymize.VerifyHash(hashParam, urlParam, salt) {
		// User error.
		log.Debug().
			Str("url", urlParam).
			Str("hash", hashParam).
			Msg("Invalid hash")
		return writeResponse(w, http.StatusUnauthorized, "invalid hash")
	}

	// Parse the url.
	target, err := url.Parse(urlParam)
	if err != nil {
		// User error.
		log.Debug().
			Str("url", urlParam).
			Msg("Invalid url")
		return writeResponse(w, http.StatusBadRequest, "invalid url")
	}

	// Get random UserAgent with corresponding Sec-Ch-Ua headers.
	ua := useragent.RandomUserAgentWithHeaders()

	// Create a new request.
	nr := &http.Request{
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

	log.Trace().
		Caller().
		Str("request", fmt.Sprint(nr)).
		Msg("Created a new request")

	// Create reverse proxy with timeout.
	rp := httputil.ReverseProxy{Director: func(r *http.Request) {}}
	rp.Transport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   timeouts.Dial,
			KeepAlive: timeouts.KeepAlive,
		}).DialContext,
		TLSHandshakeTimeout: timeouts.TLSHandshake,
	}

	// Proxy the request.
	log.Debug().
		Str("url", target.String()).
		Msg("Proxying request")
	rp.ServeHTTP(w, nr) // Use the new request.

	return nil
}
