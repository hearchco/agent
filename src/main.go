package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"

	"github.com/hearchco/agent/src/cache"
	"github.com/hearchco/agent/src/cli"
	"github.com/hearchco/agent/src/config"
	"github.com/hearchco/agent/src/logger"
	"github.com/hearchco/agent/src/profiler"
	"github.com/hearchco/agent/src/router"
)

var (
	// Release variables.
	Version   string
	Timestamp string
	GitCommit string
)

func main() {
	// Setup signal interrupt (CTRL+C).
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Parse cli flags.
	cliFlags, ver := cli.Setup(Version, Timestamp, GitCommit)

	// Configure logger.
	lgr := logger.Setup(cliFlags.Verbosity, cliFlags.Pretty)

	// Load config file.
	conf := config.New()
	conf.Load(cliFlags.ConfigPath)

	// Setup cache database.
	db, err := cache.New(ctx, conf.Server.Cache)
	if err != nil {
		log.Fatal().
			Caller().
			Err(err).
			Msg("Failed creating a new cache database")
		// ^FATAL
	}
	defer db.Close()

	// Start profiler if enabled.
	_, stopProfiler := profiler.Run(cliFlags)
	defer stopProfiler()

	// Start router.
	rw := router.New(lgr, conf, db, cliFlags.ProfilerServe, ver)
	switch conf.Server.Environment {
	case "lambda":
		rw.StartLambda()
	default:
		rw.Start(ctx)
	}

	log.Info().Msg("Program finished")
}
