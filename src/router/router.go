package router

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/hearchco/hearchco/src/cache"
	"github.com/hearchco/hearchco/src/config"
)

// it's okay to store pointer since fiber.New() returns a pointer
type RouterWrapper struct {
	app  *fiber.App
	port int
}

func New(lgr zerolog.Logger, conf config.Config, db cache.DB, serveProfiler bool) RouterWrapper {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	setupMiddlewares(app, lgr, conf.Server.FrontendUrls)
	setupRoutes(app, db, conf.Server.Cache.TTL, conf.Settings, conf.Categories)

	if serveProfiler {
		app.Use(pprof.New())
	}

	// // create new graceful engine with config port
	// rtr, err := graceful.New(gengine, graceful.WithAddr(":"+strconv.Itoa(serverConf.Port)))
	// if err != nil {
	// 	log.Error().
	// 		Err(err).
	// 		Msg("router.New(): failed creating new graceful router")
	// }

	return RouterWrapper{app: app, port: conf.Server.Port}
}

func (rw RouterWrapper) Start(ctx context.Context) {
	// shutdown on signal interrupt
	var serverShutdown sync.WaitGroup
	go func() {
		<-ctx.Done()
		log.Info().Msg("Gracefully shutting down router...")
		serverShutdown.Add(1)
		defer serverShutdown.Done()
		err := rw.app.ShutdownWithTimeout(60 * time.Second)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Router shut down failed")
		} else {
			log.Info().
				Msg("Router shut down")
		}
	}()

	// startup
	log.Info().
		Int("port", rw.port).
		Msg("Started router")

	err := rw.app.Listen(":" + strconv.Itoa(rw.port))
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed starting the router")
	}

	// wait for graceful shutdown with timeout
	serverShutdown.Wait()
}
