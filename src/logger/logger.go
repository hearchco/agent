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
	// DEBUG
	case verbosity == 1:
		log.Logger = logger.Level(zerolog.DebugLevel)
	// TRACE
	case verbosity > 1:
		log.Logger = logger.Level(zerolog.TraceLevel)
	// INFO
	default:
		log.Logger = logger.Level(zerolog.InfoLevel)
	}

	return logger
}
