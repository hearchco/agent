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
			Caller().
			Err(err).
			Msg("Failed converting string to int")
		// ^PANIC
	}
	return int64(i)
}

func convertToDurationWithoutLastChar(s string) time.Duration {
	return time.Duration(handleAtoi(s[:len(s)-1]))
}

/*
Converts the following to time.Duration:

	"1y" -> 1 year,
	"2M" -> 2 months,
	"3w" -> 3 weeks,
	"4d" -> 4 days,
	"5h" -> 5 hours,
	"6m" -> 6 minutes,
	"7s" -> 7 seconds,
	"8"-> 8 milliseconds
*/
func ConvertFromFancyTime(fancy string) time.Duration {
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

// Converts to milliseconds.
func ConvertToFancyTime(d time.Duration) string {
	return strconv.Itoa(int(d.Milliseconds()))
}
