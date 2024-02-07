package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hearchco/hearchco/src/cache"
	"github.com/hearchco/hearchco/src/cache/nocache"
	"github.com/hearchco/hearchco/src/cache/pebble"
	"github.com/hearchco/hearchco/src/cache/redis"
	"github.com/hearchco/hearchco/src/cli"
	"github.com/hearchco/hearchco/src/config"
	"github.com/hearchco/hearchco/src/logger"
	"github.com/hearchco/hearchco/src/router"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	mainTimer := time.Now()

	// setup signal interrupt (CTRL+C)
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// configure logging without file at INFO level
	_ = logger.Setup(0)

	// parse cli arguments
	cliFlags := cli.Setup()

	// start profiler
	_, stopProfiler := runProfiler(&cliFlags)
	defer stopProfiler()

	var lgr zerolog.Logger
	// configure verbosity and logging to file
	if cliFlags.LogToFile || cliFlags.Cli {
		lgr = logger.Setup(cliFlags.Verbosity, cliFlags.LogDirPath)
	} else {
		lgr = logger.Setup(cliFlags.Verbosity)
	}

	// load config file
	conf := config.New()
	conf.Load(cliFlags.DataDirPath, cliFlags.LogDirPath)

	// setup cache
	var db cache.DB
	switch conf.Server.Cache.Type {
	case "pebble":
		db = pebble.New(cliFlags.DataDirPath)
	case "redis":
		db = redis.New(ctx, conf.Server.Cache.Redis)
	default:
		db = nocache.New()
		log.Warn().Msg("Running without caching!")
	}

	// startup
	if cliFlags.Cli {
		cli.Run(cliFlags, db, conf)
	} else {
		if rw, err := router.New(conf, cliFlags.Verbosity, lgr); err != nil {
			log.Fatal().Err(err).Msg("main.main(): failed creating a router")
			// ^FATAL
		} else if err := rw.Start(ctx, db, cliFlags.ServeProfiler); err != nil {
			log.Fatal().Err(err).Msg("main.main(): failed starting the router")
			// ^FATAL
		}
	}

	// program cleanup
	db.Close()

	log.Debug().
		Int64("ms", time.Since(mainTimer).Milliseconds()).
		Msg("Program finished")
}
