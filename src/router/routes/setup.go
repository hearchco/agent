package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"

	"github.com/hearchco/agent/src/cache"
	"github.com/hearchco/agent/src/config"
)

const (
	healthzRoute            = "/healthz"
	versionzRoute           = "/versionz"
	searchWebRoute          = "/search/web"
	searchImagesRoute       = "/search/images"
	searchSuggestionsRoute  = "/search/suggestions"
	exchangeRoute           = "/exchange"
	exchangeCurrenciesRoute = "/exchange/currencies"
	imageProxyRoute         = "/imageproxy"
)

func Setup(mux *chi.Mux, ver string, db cache.DB, conf config.Config) {
	// Health check
	mux.Get(healthzRoute, func(w http.ResponseWriter, r *http.Request) {
		err := writeResponse(w, http.StatusOK, "OK")
		if err != nil {
			log.Error().
				Err(err).
				Str("path", r.URL.Path).
				Str("method", r.Method).
				Msg("Failed to send response")
		}
	})

	// Version
	mux.Get(versionzRoute, func(w http.ResponseWriter, r *http.Request) {
		err := writeResponse(w, http.StatusOK, ver)
		if err != nil {
			log.Error().
				Err(err).
				Str("path", r.URL.Path).
				Str("method", r.Method).
				Msg("Failed to send response")
		}
	})

	// Web search
	muxGetPost(mux, searchWebRoute, func(w http.ResponseWriter, r *http.Request) {
		err := routeSearchWeb(w, r, ver, conf.Engines.NoWeb, conf.Server.ImageProxy.SecretKey)
		if err != nil {
			log.Error().
				Err(err).
				Str("path", r.URL.Path).
				Str("method", r.Method).
				Msg("Failed to send response")
		}
	})

	// Images search
	muxGetPost(mux, searchImagesRoute, func(w http.ResponseWriter, r *http.Request) {
		err := routeSearchImages(w, r, ver, conf.Engines.NoImages, conf.Server.ImageProxy.SecretKey)
		if err != nil {
			log.Error().
				Err(err).
				Str("path", r.URL.Path).
				Str("method", r.Method).
				Msg("Failed to send response")
		}
	})

	// Suggestions search
	muxGetPost(mux, searchSuggestionsRoute, func(w http.ResponseWriter, r *http.Request) {
		err := routeSearchSuggestions(w, r, ver, conf.Engines.NoSuggestions)
		if err != nil {
			log.Error().
				Err(err).
				Str("path", r.URL.Path).
				Str("method", r.Method).
				Msg("Failed to send response")
		}
	})

	// Exchange
	muxGetPost(mux, exchangeRoute, func(w http.ResponseWriter, r *http.Request) {
		err := routeExchange(w, r, ver, conf.Exchange, db, conf.Server.Cache.TTL.Currencies)
		if err != nil {
			log.Error().
				Err(err).
				Str("path", r.URL.Path).
				Str("method", r.Method).
				Msg("Failed to send response")
		}
	})

	// Exchange currencies
	muxGetPost(mux, exchangeCurrenciesRoute, func(w http.ResponseWriter, r *http.Request) {
		err := routeCurrencies(w, ver, conf.Exchange, db, conf.Server.Cache.TTL.Currencies)
		if err != nil {
			log.Error().
				Err(err).
				Str("path", r.URL.Path).
				Str("method", r.Method).
				Msg("Failed to send response")
		}
	})

	// Image proxy
	mux.Get(imageProxyRoute, func(w http.ResponseWriter, r *http.Request) {
		err := routeImageProxy(w, r, conf.Server.ImageProxy.SecretKey, conf.Server.ImageProxy.Timeout)
		if err != nil {
			log.Error().
				Err(err).
				Str("path", r.URL.Path).
				Str("method", r.Method).
				Msg("Failed to send response")
		}
	})
}

func muxGetPost(mux *chi.Mux, pattern string, handler http.HandlerFunc) {
	mux.Get(pattern, handler)
	mux.Post(pattern, handler)
}
