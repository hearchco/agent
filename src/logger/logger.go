package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Setup(verbosity int8, pretty bool) zerolog.Logger {
	// Setup logger
	var logger zerolog.Logger
	// if pretty use console writer
	if pretty {
		logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.Stamp})
	} else {
		logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	}

	// Setup verbosity
	switch {
	// TRACE
	case verbosity > 1:
		logger = logger.Level(zerolog.TraceLevel)
	// DEBUG
	case verbosity == 1:
		logger = logger.Level(zerolog.DebugLevel)
	// INFO
	default:
		logger = logger.Level(zerolog.InfoLevel)
	}

	// set the logger to global and return it
	log.Logger = logger
	return logger
}
