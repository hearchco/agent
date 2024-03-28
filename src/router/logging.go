package router

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

func ignoredPath(path string, skipPaths []string) bool {
	for _, p := range skipPaths {
		if p == path {
			return true
		}
	}
	return false
}

func zerologMiddleware(lgr zerolog.Logger, skipPaths []string) [](func(http.Handler) http.Handler) {
	newHandler := hlog.NewHandler(lgr)
	fieldsHandler := hlog.AccessHandler(func(r *http.Request, status int, size int, duration time.Duration) {
		// skip logging for ignored paths
		if ignoredPath(r.URL.Path, skipPaths) {
			return
		}

		// get logger from context
		lgr := hlog.FromRequest(r)

		// decide on log level
		event := lgr.Info()
		if status >= 500 {
			event = lgr.Error()
		} else if status >= 400 {
			event = lgr.Warn()
		}

		// log
		event.
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status", status).
			Dur("duration", duration).
			Str("ip", r.RemoteAddr).
			Msg("Request")
	})

	return [](func(http.Handler) http.Handler){
		newHandler,
		fieldsHandler,
	}
}
