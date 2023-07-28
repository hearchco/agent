package search

import (
	"errors"

	"golang.org/x/time/rate"
)

// ErrRateLimited indicates that you have been detected of scraping and temporarily blocked.
// The duration of the block is unspecified.
var ErrRateLimited = errors.New("ratelimited")

// RateLimit sets a global limit to how many requests can be made in a given time interval.
// The default is unlimited (but obviously you will get blocked temporarily if you do too many
// calls too quickly).
//
// See: https://godoc.org/golang.org/x/time/rate#NewLimiter
func CreateRateLimit(limit float64) *rate.Limiter {
	if limit == 0 {
		return rate.NewLimiter(rate.Inf, 0)
	} else {
		return rate.NewLimiter(rate.Limit(limit), 0)
	}
}
