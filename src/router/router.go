package router

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/hearchco/hearchco/src/cache"
	"github.com/hearchco/hearchco/src/config"
)

// it's okay to store pointer since chi.NewRouter() returns a pointer
type RouterWrapper struct {
	mux  *chi.Mux
	port int
}

func New(lgr zerolog.Logger, conf config.Config, db cache.DB, serveProfiler bool) RouterWrapper {
	mux := chi.NewRouter()

	setupMiddlewares(mux, lgr, conf.Server.FrontendUrls, serveProfiler)
	setupRoutes(mux, db, conf)

	return RouterWrapper{mux: mux, port: conf.Server.Port}
}

func (rw RouterWrapper) Start(ctx context.Context) {
	// create server
	srv := http.Server{
		Addr:    ":" + strconv.Itoa(rw.port),
		Handler: rw.mux,
	}

	log.Info().
		Int("port", rw.port).
		Msg("Starting server")

	// shut down server gracefully on context cancellation
	go func() {
		<-ctx.Done()
		log.Info().Msg("Shutting down server")

		// create a context with timeout of 5 seconds
		timeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// shutdown gracefully
		// after timeout is reached, server will be shut down forcefully
		err := srv.Shutdown(timeout)
		if err != nil {
			log.Error().
				Caller().
				Err(err).
				Msg("Server shut down failed")
		} else {
			log.Info().
				Msg("Server shut down")
		}
	}()

	// start server
	err := srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal().
			Caller().
			Err(err).
			Msg("Failed to start server")
	}
}
