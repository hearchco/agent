package moretime

import (
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

func handleAtoi(s string) int64 {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Panic().
			Err(err).
			Msg("failed converting string to int")
		// ^PANIC
	}
	return int64(i)
}

func convertToDurationWithoutLastChar(s string) time.Duration {
	return time.Duration(handleAtoi(s[:len(s)-1]))
}

func ConvertFancyTime(fancy string) time.Duration {
	switch fancy[len(fancy)-1] {
	case 'y':
		return convertToDurationWithoutLastChar(fancy) * Year
	case 'M':
		return convertToDurationWithoutLastChar(fancy) * Month
	case 'w':
		return convertToDurationWithoutLastChar(fancy) * Week
	case 'd':
		return convertToDurationWithoutLastChar(fancy) * Day
	case 'h':
		return convertToDurationWithoutLastChar(fancy) * time.Hour
	case 'm':
		return convertToDurationWithoutLastChar(fancy) * time.Minute
	case 's':
		return convertToDurationWithoutLastChar(fancy) * time.Second
	default:
		return time.Duration(handleAtoi(fancy)) * time.Millisecond
	}
}

func ConvertToFancyTime(d time.Duration) string {
	return strconv.Itoa(int(d.Milliseconds()))
}
