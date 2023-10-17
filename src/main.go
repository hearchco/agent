package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/tminaorg/brzaguza/src/cache"
	"github.com/tminaorg/brzaguza/src/cache/nocache"
	"github.com/tminaorg/brzaguza/src/cache/pebble"
	"github.com/tminaorg/brzaguza/src/cache/redis"
	"github.com/tminaorg/brzaguza/src/cli"
	"github.com/tminaorg/brzaguza/src/config"
	"github.com/tminaorg/brzaguza/src/logger"
	"github.com/tminaorg/brzaguza/src/router"
)

func main() {
	mainTimer := time.Now()

	// parse cli arguments
	cliFlags := cli.Setup()

	_, stopProfiler := runProfiler(&cliFlags)
	defer stopProfiler()

	// signal interrupt (CTRL+C)
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// configure logging
	logger.Setup(cliFlags.Log, cliFlags.Verbosity)

	// load config file
	conf := config.New()
	conf.Load(cliFlags.Data, cliFlags.Log)

	// cache database
	var db cache.DB
	switch conf.Server.Cache.Type {
	case "pebble":
		db = pebble.New(cliFlags.Data)
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
			log.Error().Err(err).Msg("Failed creating a router")
		} else {
			rw.Start(ctx, db, cliFlags.ServeProfiler)
		}
	}

	// program cleanup
	db.Close()

	log.Debug().Msgf("Program finished in %vms", time.Since(mainTimer).Milliseconds())
}
