package cache

import (
	"time"

	"github.com/hearchco/hearchco/src/search/result"
)

type DB interface {
	Close()
	Set(k string, v any, ttl ...time.Duration) error
	SetResults(query string, category string, results []result.Result, ttl ...time.Duration) error
	Get(k string, o any) error
	GetResults(query string, category string) ([]result.Result, error)
	GetTTL(k string) (time.Duration, error)
	GetResultsTTL(query string, category string) (time.Duration, error)
}
