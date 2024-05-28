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

// converts 1y to 1 year
// converts 2M to 2 month
// converts 3w to 3 week
// converts 4d to 4 day
// converts 5h to 5 hour
// converts 6m to 6 minute
// converts 7s to 7 second
// converts 8 to 8 millisecond
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

// converts to milliseconds
func ConvertToFancyTime(d time.Duration) string {
	return strconv.Itoa(int(d.Milliseconds()))
}
