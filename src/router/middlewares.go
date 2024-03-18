package router

import (
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

	// use trailing slash middleware
	mux.Use(middleware.StripSlashes)

	// use gzip middleware
	mux.Use(middleware.Compress(5))

	// use CORS middleware
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   frontendUrls,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Content-Encoding", "Accept-Encoding"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	log.Debug().
		Strs("url", frontendUrls).
		Msg("Using CORS")

	// use profiler middleware if enabled
	if serveProfiler {
		mux.Mount("/debug", middleware.Profiler())
	}
}
