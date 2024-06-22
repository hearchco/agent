package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"

	"github.com/hearchco/agent/src/cache"
	"github.com/hearchco/agent/src/config"
)

func Setup(mux *chi.Mux, ver string, db cache.DB, conf config.Config) {
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

	mux.Get("/search", func(w http.ResponseWriter, r *http.Request) {
		err := routeSearch(w, r, ver, conf.Categories, conf.Server.Cache.TTL, db, conf.Server.ImageProxy.Salt)
		if err != nil {
			log.Error().
				Err(err).
				Str("path", r.URL.Path).
				Str("method", r.Method).
				Msg("Failed to send response")
		}
	})

	mux.Post("/suggestions", func(w http.ResponseWriter, r *http.Request) {
		err := routeSuggest(w, r)
		if err != nil {
			log.Error().
				Err(err).
				Str("path", r.URL.Path).
				Str("method", r.Method).
				Msg("Failed to send response")
		}
	})

	mux.Get("/suggestions", func(w http.ResponseWriter, r *http.Request) {
		err := routeSuggest(w, r)
		if err != nil {
			log.Error().
				Err(err).
				Str("path", r.URL.Path).
				Str("method", r.Method).
				Msg("Failed to send response")
		}
	})

	mux.Post("/search", func(w http.ResponseWriter, r *http.Request) {
		err := routeSearch(w, r, ver, conf.Categories, conf.Server.Cache.TTL, db, conf.Server.ImageProxy.Salt)
		if err != nil {
			log.Error().
				Err(err).
				Str("path", r.URL.Path).
				Str("method", r.Method).
				Msg("Failed to send response")
		}
	})

	mux.Get("/proxy", func(w http.ResponseWriter, r *http.Request) {
		err := routeProxy(w, r, conf.Server.ImageProxy.Salt, conf.Server.ImageProxy.Timeouts)
		if err != nil {
			log.Error().
				Err(err).
				Str("path", r.URL.Path).
				Str("method", r.Method).
				Msg("Failed to send response")
		}
	})
}
