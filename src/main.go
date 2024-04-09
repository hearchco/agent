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
	"github.com/rs/zerolog"
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
	_, stopProfiler := runProfiler(cliFlags)
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
	db, err := cache.New(ctx, cliFlags.DataDirPath, conf.Server.Cache)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("main.main(): failed creating a new db")
		// ^FATAL
	}

	// startup
	if cliFlags.Cli {
		cli.Run(cliFlags, db, conf)
	} else {
		rw := router.New(lgr, conf, db, cliFlags.ServeProfiler)
		rw.Start(ctx)
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
