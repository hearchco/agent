package middlewares

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Setup(mux *chi.Mux, lgr zerolog.Logger, frontendUrls []string, serveProfiler bool) {
	// Use custom zerolog middleware.
	// TODO: Make skipped paths configurable.
	skipPaths := []string{"/healthz", "/versionz"}
	mux.Use(zerologMiddleware(lgr, skipPaths)...)

	// Use recovery middleware.
	mux.Use(middleware.Recoverer)

	// Use compression middleware, except for image proxy since the response is copied over.
	mux.Use(middleware.Maybe(compress(3), func(r *http.Request) bool {
		return !strings.HasPrefix(r.URL.Path, "/proxy")
	}))

	// Use CORS middleware.
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins: frontendUrls,
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{
			"Accept",
			"Accept-Encoding",
			"Accept-Language",
			"Access-Control-Request-Headers",
			"Access-Control-Request-Method",
			"Origin",
		},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	log.Debug().
		Strs("url", frontendUrls).
		Msg("Using CORS")

	// Use strip slashes middleware, except for pprof.
	mux.Use(middleware.Maybe(middleware.StripSlashes, func(r *http.Request) bool {
		return !strings.HasPrefix(r.URL.Path, "/debug")
	}))

	// Use pprof router if profiling is enabled.
	if serveProfiler {
		mux.Mount("/debug", middleware.Profiler())
	}
}
