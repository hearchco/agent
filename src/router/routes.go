package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"

	"github.com/hearchco/hearchco/src/cache"
	"github.com/hearchco/hearchco/src/config"
)

func setupRoutes(mux *chi.Mux, db cache.DB, conf config.Config) {
	mux.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		err := writeResponse(w, http.StatusOK, "OK")
		if err != nil {
			log.Error().
				Err(err).
				Str("path", r.URL.Path).
				Str("method", r.Method).
				Msg("Failed to healthz")
		}
	})

	mux.Get("/search", func(w http.ResponseWriter, r *http.Request) {
		err := Search(w, r, db, conf.Server.Cache.TTL, conf.Settings, conf.Categories, conf.Server.Proxy.Salt)
		if err != nil {
			log.Error().
				Err(err).
				Str("path", r.URL.Path).
				Str("method", r.Method).
				Msg("Failed to search")
		}
	})

	mux.Post("/search", func(w http.ResponseWriter, r *http.Request) {
		err := Search(w, r, db, conf.Server.Cache.TTL, conf.Settings, conf.Categories, conf.Server.Proxy.Salt)
		if err != nil {
			log.Error().
				Err(err).
				Str("path", r.URL.Path).
				Str("method", r.Method).
				Msg("Failed to search")
		}
	})

	mux.Get("/proxy", func(w http.ResponseWriter, r *http.Request) {
		err := Proxy(w, r, conf.Server.Proxy.Salt, conf.Server.Proxy.Timeouts)
		if err != nil {
			log.Error().
				Err(err).
				Str("path", r.URL.Path).
				Str("method", r.Method).
				Msg("Failed to proxy")
		}
	})
}
