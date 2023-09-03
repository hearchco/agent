package logger

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func CalculateDatetime() string {
	year, month, day := time.Now().Date()
	datetime := fmt.Sprintf("%v%v%v", year, month, day)
	if month < 10 {
		if day < 10 {
			datetime = fmt.Sprintf("%v0%v0%v", year, month, day)
		} else {
			datetime = fmt.Sprintf("%v0%v%v", year, month, day)
		}
	} else {
		if day < 10 {
			datetime = fmt.Sprintf("%v%v0%v", year, month, day)
		}
	}
	return datetime
}

func Setup(path string, name string, verbosity int) {
	// Check if path ends with "/" and add it otherwise
	if path[len(path)-1] != '/' {
		path = path + "/"
	}

	// Generate logfile name
	datetime := CalculateDatetime()
	fullpath := fmt.Sprintf("%v%v_%v.log", path, name, datetime)

	// Setup logger
	logger := log.Output(io.MultiWriter(zerolog.ConsoleWriter{
		TimeFormat: time.Stamp,
		Out:        os.Stderr,
	}, zerolog.ConsoleWriter{
		TimeFormat: time.Stamp,
		Out: &lumberjack.Logger{
			Filename:   fullpath,
			MaxSize:    5,
			MaxAge:     14,
			MaxBackups: 5,
		},
		NoColor: true,
	}))

	// Setup verbosity
	switch {
	case verbosity == 1:
		log.Logger = logger.Level(zerolog.DebugLevel)
	case verbosity > 1:
		log.Logger = logger.Level(zerolog.TraceLevel)
	default:
		log.Logger = logger.Level(zerolog.InfoLevel)
	}
}
