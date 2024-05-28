package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hearchco/hearchco/src/cache"
	"github.com/hearchco/hearchco/src/cli"
	"github.com/hearchco/hearchco/src/config"
	"github.com/hearchco/hearchco/src/logger"
	"github.com/hearchco/hearchco/src/router"
	"github.com/rs/zerolog/log"
)

func main() {
	mainTimer := time.Now()

	// setup signal interrupt (CTRL+C)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// parse cli arguments
	cliFlags := cli.Setup()

	// configure logger
	lgr := logger.Setup(cliFlags.Verbosity, cliFlags.Pretty)

	// start profiler
	_, stopProfiler := runProfiler(cliFlags)
	defer stopProfiler()

	// load config file
	conf := config.New()
	conf.Load(cliFlags.DataDirPath)

	// setup cache
	db, err := cache.New(ctx, cliFlags.DataDirPath, conf.Server.Cache)
	if err != nil {
		log.Fatal().
			Caller().
			Err(err).
			Msg("Failed creating a new db")
		// ^FATAL
	}

	// startup
	if cliFlags.Cli {
		cli.Run(cliFlags, db, conf)
	} else {
		rw := router.New(lgr, conf, db, cliFlags.ServeProfiler)
		switch conf.Server.Environment {
		case "lambda":
			rw.StartLambda()
		default:
			rw.Start(ctx)
		}
	}

	// program cleanup
	db.Close()

	if cliFlags.Cli {
		log.Debug().
			Dur("duration", time.Since(mainTimer)).
			Msg("Program finished")
	} else {
		log.Info().
			Dur("duration", time.Since(mainTimer)).
			Msg("Program finished")
	}
}
