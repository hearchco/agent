package router

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/hearchco/hearchco/src/anonymize"
	"github.com/hearchco/hearchco/src/search/useragent"
	"github.com/rs/zerolog/log"
)

func Proxy(w http.ResponseWriter, r *http.Request, salt string, timeout time.Duration) error {
	err := r.ParseForm()
	if err != nil {
		// server error
		writeResponse(w, http.StatusInternalServerError, "failed to parse form")
		return err
	}

	params := r.Form

	urlParam := getParamOrDefault(params, "url")
	hashParam := getParamOrDefault(params, "hash")

	if urlParam == "" || hashParam == "" {
		// user error
		writeResponse(w, http.StatusBadRequest, "url and hash are required")
		return nil
	}

	// check if hash is valid
	if !anonymize.CheckHash(hashParam, urlParam, salt) {
		// user error
		writeResponse(w, http.StatusUnauthorized, "invalid hash")
		log.Debug().
			Str("url", urlParam).
			Str("hash", hashParam).
			Msg("Invalid hash")
		return nil
	}

	// parse the url
	target, err := url.Parse(urlParam)
	if err != nil {
		// user error
		writeResponse(w, http.StatusBadRequest, "invalid url")
		log.Debug().
			Str("url", urlParam).
			Msg("Invalid url")
		return nil
	}

	// create new request
	nr := &http.Request{
		URL: target,
		Header: map[string][]string{
			"User-Agent":      {useragent.RandomUserAgent()},
			"Accept":          {"image/avif", "image/webp", "*/*"},
			"Accept-Encoding": {"gzip", "deflate", "br"},
			"Sec-GPC":         {"1"}, // don't share info with 3rd parties
			"DNT":             {"1"}, // do not track
		},
	}

	log.Trace().
		Str("request", fmt.Sprint(nr)).
		Msg("Created new request")

	// TODO: implement timeout
	// proxy the request
	rp := httputil.ReverseProxy{Director: func(r *http.Request) {}}
	log.Debug().
		Str("url", target.String()).
		Msg("Proxying request")
	rp.ServeHTTP(w, nr) // use new request

	return nil
}
