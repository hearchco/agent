package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Setup(verbosity int8, pretty bool) zerolog.Logger {
	// Setup logger
	var lgr zerolog.Logger
	if pretty {
		lgr = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.Stamp})
	} else {
		lgr = zerolog.New(os.Stderr).With().Timestamp().Logger()
	}

	// Setup verbosity
	switch {
	// TRACE
	case verbosity > 1:
		lgr = lgr.With().Caller().Logger().Level(zerolog.TraceLevel)
	// DEBUG
	case verbosity == 1:
		lgr = lgr.Level(zerolog.DebugLevel)
	// INFO
	default:
		lgr = lgr.Level(zerolog.InfoLevel)
	}

	// Set the logger to global and return it
	log.Logger = lgr
	return lgr
}
