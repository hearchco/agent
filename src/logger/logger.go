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
	year, monthS, day := time.Now().Date()
	month := int(monthS)

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

func Setup(path string, verbosity int) {
	// Generate logfile name
	datetime := CalculateDatetime()
	filepath := fmt.Sprintf("%v/log/brzaguza_%v.log", path, datetime)

	// Setup logger
	logger := log.Output(io.MultiWriter(zerolog.ConsoleWriter{
		TimeFormat: time.Stamp,
		Out:        os.Stderr,
	}, zerolog.ConsoleWriter{
		TimeFormat: time.Stamp,
		Out: &lumberjack.Logger{
			Filename:   filepath,
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
