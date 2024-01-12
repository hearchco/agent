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

func Setup(verbosity int8, logDirPath ...string) {
	logWriters := []io.Writer{zerolog.ConsoleWriter{
		TimeFormat: time.Stamp,
		Out:        os.Stderr,
	}}

	// Generate logfile name and ConsoleWriter to file
	if logDirPath[0] != "" {
		logFilePath := path.Join(logDirPath[0], fmt.Sprintf("hearchco_%v.log", time.Now().Format("20060102")))
		logWriters = append(logWriters, zerolog.ConsoleWriter{
			TimeFormat: time.Stamp,
			Out: &lumberjack.Logger{
				Filename:   logFilePath,
				MaxSize:    5,
				MaxAge:     14,
				MaxBackups: 5,
			},
			NoColor: true,
		})
	}

	// Setup logger
	logger := log.Output(io.MultiWriter(logWriters...))

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
