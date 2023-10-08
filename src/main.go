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
	"github.com/tminaorg/brzaguza/src/climode"
	"github.com/tminaorg/brzaguza/src/config"
	"github.com/tminaorg/brzaguza/src/logger"
	"github.com/tminaorg/brzaguza/src/router"
)

func main() {
	mainTimer := time.Now()

	// parse cli arguments
	setupCli()

	// configure logging
	logger.Setup(cli.Log, cli.Verbosity)

	// load config file
	conf := config.New()
	conf.Load(cli.Config, cli.Log)

	// signal interrupt (CTRL+C)
	ctx, stopCtx := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// profiler
	_, stopProfiler := runProfiler() // not used currently

	// cache database
	var db cache.DB
	switch conf.Server.Cache.Type {
	case "pebble":
		db = pebble.New(cli.Config)
	case "redis":
		db = redis.New(ctx, conf.Server.Cache.Redis)
	default:
		db = nocache.New()
		log.Warn().Msg("Running without caching!")
	}

	// startup
	if cli.Cli {
		climode.Run(cli.Query, cli.MaxPages, cli.Visit, cli.Silent, db, conf)
	} else {
		if rw, err := router.New(conf, cli.Verbosity); err != nil {
			log.Error().Msgf("Failed creating a router: %v", err)
		} else {
			rw.Start(ctx, db, cli.ServeProfiler)
		}
	}

	// program cleanup
	db.Close()
	stopCtx()
	stopProfiler()

	log.Debug().Msgf("Program finished in %vms", time.Since(mainTimer).Milliseconds())
}
