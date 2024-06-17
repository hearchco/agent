package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Setup(verbosity int8, pretty bool) zerolog.Logger {
	// Setup logger.
	var l zerolog.Logger
	if pretty {
		// This is much slower to print.
		l = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.Stamp})
	} else {
		l = zerolog.New(os.Stderr).With().Timestamp().Logger()
	}

	// Setup verbosity.
	switch {
	case verbosity > 1: // TRACE
		l = l.With().Caller().Logger().Level(zerolog.TraceLevel)
	case verbosity == 1: // DEBUG
		l = l.Level(zerolog.DebugLevel)
	default: // INFO
		l = l.Level(zerolog.InfoLevel)
	}

	// Set the logger to be global.
	log.Logger = l
	return l
}
