package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hearchco/hearchco/src/cache"
	"github.com/hearchco/hearchco/src/cache/badger"
	"github.com/hearchco/hearchco/src/cache/nocache"
	"github.com/hearchco/hearchco/src/cache/redis"
	"github.com/hearchco/hearchco/src/cli"
	"github.com/hearchco/hearchco/src/config"
	"github.com/hearchco/hearchco/src/logger"
	"github.com/hearchco/hearchco/src/router"
	"github.com/rs/zerolog/log"
)

func main() {
	mainTimer := time.Now()

	// setup signal interrupt (CTRL+C)
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// configure logging without file at INFO level
	logger.Setup(0)

	// parse cli arguments
	cliFlags := cli.Setup()

	// start profiler
	_, stopProfiler := runProfiler(&cliFlags)
	defer stopProfiler()

	// configure verbosity and logging to file
	if cliFlags.LogToFile || cliFlags.Cli {
		logger.Setup(cliFlags.Verbosity, cliFlags.LogDirPath)
	} else {
		logger.Setup(cliFlags.Verbosity)
	}

	// load config file
	conf := config.New()
	conf.Load(cliFlags.DataDirPath, cliFlags.LogDirPath)

	// setup cache
	var db cache.DB
	switch conf.Server.Cache.Type {
	case "badger":
		db = badger.New(cliFlags.DataDirPath, conf.Server.Cache.Badger.Persist)
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
		if rw, err := router.New(conf, cliFlags.Verbosity); err != nil {
			log.Fatal().Err(err).Msg("main.main(): failed creating a router")
			// ^FATAL
		} else {
			rw.Start(ctx, db, cliFlags.ServeProfiler)
		}
	}

	// program cleanup
	db.Close()

	log.Debug().
		Int64("ms", time.Since(mainTimer).Milliseconds()).
		Msg("Program finished")
}
