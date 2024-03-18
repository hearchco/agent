package router

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/hearchco/hearchco/src/anonymize"
	"github.com/rs/zerolog/log"
)

func Proxy(w http.ResponseWriter, r *http.Request, salt string, timeout time.Duration) error {
	r.ParseForm()
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
		return nil
	}

	// parse the url
	target, err := url.Parse(urlParam)
	if err != nil {
		// user error
		writeResponse(w, http.StatusBadRequest, "invalid url")
		return nil
	}

	log.Debug().
		Str("url", target.String()).
		Msg("Proxying request")

	// TODO: implement timeout
	// proxy the request
	rp := httputil.ReverseProxy{Director: func(r *http.Request) {
		r.URL.Scheme = target.Scheme
		r.URL.Host = target.Host
		r.URL.Path = target.Path
		r.Host = target.Host
	}}
	rp.ServeHTTP(w, r)

	return nil
}
