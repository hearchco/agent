package logger

import (
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func DateString() string {
	return time.Now().Format("20060102")
}

func Setup(logDirPath string, verbosity int8) {
	// Generate logfile name
	datetime := DateString()
	logFilePath := path.Join(logDirPath, fmt.Sprintf("hearchco_%v.log", datetime))

	// Setup logger
	logger := log.Output(io.MultiWriter(zerolog.ConsoleWriter{
		TimeFormat: time.Stamp,
		Out:        os.Stderr,
	}, zerolog.ConsoleWriter{
		TimeFormat: time.Stamp,
		Out: &lumberjack.Logger{
			Filename:   logFilePath,
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
