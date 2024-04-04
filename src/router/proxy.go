package router

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/hearchco/hearchco/src/anonymize"
	"github.com/hearchco/hearchco/src/config"
	"github.com/hearchco/hearchco/src/search/useragent"
	"github.com/rs/zerolog/log"
)

func Proxy(w http.ResponseWriter, r *http.Request, salt string, timeouts config.ProxyTimeouts) error {
	err := r.ParseForm()
	if err != nil {
		// server error
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
		// user error
		return writeResponse(w, http.StatusBadRequest, "url and hash are required")
	}

	// check if hash is valid
	if !anonymize.CheckHash(hashParam, urlParam, salt) {
		// user error
		log.Debug().
			Str("url", urlParam).
			Str("hash", hashParam).
			Msg("Invalid hash")
		return writeResponse(w, http.StatusUnauthorized, "invalid hash")
	}

	// parse the url
	target, err := url.Parse(urlParam)
	if err != nil {
		// user error
		log.Debug().
			Str("url", urlParam).
			Msg("Invalid url")
		return writeResponse(w, http.StatusBadRequest, "invalid url")
	}

	// create new request
	nr := &http.Request{
		Method: http.MethodGet,
		URL:    target,
		Host:   target.Host,
		// RemoteAddr: "127.0.0.1", // TODO: implement server IP getting (should be cached)
		RequestURI: target.RequestURI(),
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header: map[string][]string{
			"Accept":             {"image/avif", "image/webp", "*/*"},
			"Accept-Encoding":    {"gzip", "deflate", "br"},
			"Accept-Language":    {"en-US,en;q=0.5"},
			"Connection":         {"keep-alive"},
			"DNT":                {"1"}, // do not track
			"sec-ch-ua":          {`"Google Chrome";v="118", "Chromium";v="118", "Not=A?Brand";v="24"`},
			"sec-ch-ua-mobile":   {"?0"},
			"sec-ch-ua-platform": {"\"Windows\""},
			"Sec-Fetch-Dest":     {"image"},
			"Sec-Fetch-Mode":     {"no-cors"},
			"Sec-Fetch-Site":     {"same-site"},
			"Sec-GPC":            {"1"}, // don't share info with 3rd parties
			"TE":                 {"trailers"},
			"User-Agent":         {useragent.RandomUserAgent()},
		},
	}

	log.Trace().
		Str("request", fmt.Sprint(nr)).
		Msg("Created new request")

	// create reverse proxy with timeout
	rp := httputil.ReverseProxy{Director: func(r *http.Request) {}}
	rp.Transport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   timeouts.Dial,
			KeepAlive: timeouts.KeepAlive,
		}).DialContext,
		TLSHandshakeTimeout: timeouts.TLSHandshake,
	}

	// proxy the request
	log.Debug().
		Str("url", target.String()).
		Msg("Proxying request")
	rp.ServeHTTP(w, nr) // use new request

	return nil
}
