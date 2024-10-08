package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"

	"github.com/hearchco/agent/src/cache"
	"github.com/hearchco/agent/src/config"
	"github.com/hearchco/agent/src/search/category"
)

func Setup(mux *chi.Mux, ver string, db cache.DB, conf config.Config) {
	// /healthz
	mux.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		err := writeResponse(w, http.StatusOK, "OK")
		if err != nil {
			log.Error().
				Err(err).
				Str("path", r.URL.Path).
				Str("method", r.Method).
				Msg("Failed to send response")
		}
	})

	// /versionz
	mux.Get("/versionz", func(w http.ResponseWriter, r *http.Request) {
		err := writeResponse(w, http.StatusOK, ver)
		if err != nil {
			log.Error().
				Err(err).
				Str("path", r.URL.Path).
				Str("method", r.Method).
				Msg("Failed to send response")
		}
	})

	// /search
	mux.Get("/search", func(w http.ResponseWriter, r *http.Request) {
		err := routeSearch(w, r, ver, conf.Categories, conf.Server.ImageProxy.SecretKey)
		if err != nil {
			log.Error().
				Err(err).
				Str("path", r.URL.Path).
				Str("method", r.Method).
				Msg("Failed to send response")
		}
	})
	mux.Post("/search", func(w http.ResponseWriter, r *http.Request) {
		err := routeSearch(w, r, ver, conf.Categories, conf.Server.ImageProxy.SecretKey)
		if err != nil {
			log.Error().
				Err(err).
				Str("path", r.URL.Path).
				Str("method", r.Method).
				Msg("Failed to send response")
		}
	})

	// /suggestions
	mux.Get("/suggestions", func(w http.ResponseWriter, r *http.Request) {
		err := routeSuggest(w, r, ver, conf.Categories[category.SUGGESTIONS])
		if err != nil {
			log.Error().
				Err(err).
				Str("path", r.URL.Path).
				Str("method", r.Method).
				Msg("Failed to send response")
		}
	})
	mux.Post("/suggestions", func(w http.ResponseWriter, r *http.Request) {
		err := routeSuggest(w, r, ver, conf.Categories[category.SUGGESTIONS])
		if err != nil {
			log.Error().
				Err(err).
				Str("path", r.URL.Path).
				Str("method", r.Method).
				Msg("Failed to send response")
		}
	})

	// /exchange
	mux.Get("/exchange", func(w http.ResponseWriter, r *http.Request) {
		err := routeExchange(w, r, ver, conf.Exchange, db, conf.Server.Cache.TTL.Currencies)
		if err != nil {
			log.Error().
				Err(err).
				Str("path", r.URL.Path).
				Str("method", r.Method).
				Msg("Failed to send response")
		}
	})
	mux.Post("/exchange", func(w http.ResponseWriter, r *http.Request) {
		err := routeExchange(w, r, ver, conf.Exchange, db, conf.Server.Cache.TTL.Currencies)
		if err != nil {
			log.Error().
				Err(err).
				Str("path", r.URL.Path).
				Str("method", r.Method).
				Msg("Failed to send response")
		}
	})

	// /currencies
	mux.Get("/currencies", func(w http.ResponseWriter, r *http.Request) {
		err := routeCurrencies(w, ver, conf.Exchange, db, conf.Server.Cache.TTL.Currencies)
		if err != nil {
			log.Error().
				Err(err).
				Str("path", r.URL.Path).
				Str("method", r.Method).
				Msg("Failed to send response")
		}
	})
	mux.Post("/currencies", func(w http.ResponseWriter, r *http.Request) {
		err := routeCurrencies(w, ver, conf.Exchange, db, conf.Server.Cache.TTL.Currencies)
		if err != nil {
			log.Error().
				Err(err).
				Str("path", r.URL.Path).
				Str("method", r.Method).
				Msg("Failed to send response")
		}
	})

	// /proxy
	mux.Get("/proxy", func(w http.ResponseWriter, r *http.Request) {
		err := routeProxy(w, r, conf.Server.ImageProxy.SecretKey, conf.Server.ImageProxy.Timeout)
		if err != nil {
			log.Error().
				Err(err).
				Str("path", r.URL.Path).
				Str("method", r.Method).
				Msg("Failed to send response")
		}
	})
}
