package cache

import (
	"time"

	"github.com/hearchco/hearchco/src/search/result"
)

type DB interface {
	Close() error
	GetResults(query string) ([]result.Result, error)
	GetImageResults(query string, salt string) ([]result.Result, error)
	SetResults(query string, results []result.Result) error
	SetImageResults(query string, results []result.Result) error
	GetAge(query string) (time.Duration, error)
}
