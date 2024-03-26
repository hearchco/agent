package router

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func setupMiddlewares(mux *chi.Mux, lgr zerolog.Logger, frontendUrls []string, serveProfiler bool) {
	// use custom zerolog middleware
	// TODO: make skipped paths configurable
	skipPaths := []string{"/healthz"}
	mux.Use(zerologMiddleware(lgr, skipPaths)...)

	// use recovery middleware
	mux.Use(middleware.Recoverer)

	// use compression middleware
	mux.Use(compress(5)...)

	// use CORS middleware
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   frontendUrls,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Accept-Encoding"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	log.Debug().
		Strs("url", frontendUrls).
		Msg("Using CORS")

	// use strip slashes middleware except for pprof
	mux.Use(middleware.Maybe(middleware.StripSlashes, func(r *http.Request) bool {
		return !strings.HasPrefix(r.URL.Path, "/debug")
	}))

	// use pprof router if enabled
	if !serveProfiler {
		mux.Mount("/debug", middleware.Profiler())
	}
}
