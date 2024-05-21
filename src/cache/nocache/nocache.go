package nocache

import (
	"fmt"
	"time"

	"github.com/hearchco/hearchco/src/search/result"
)

type DB struct{}

func New() (DB, error) { return DB{}, nil }

func (db DB) Close() {}

func (db DB) Set(k string, v interface{}, ttl ...time.Duration) error { return nil }

func (db DB) SetResults(query string, category string, results []result.Result, ttl ...time.Duration) error {
	return db.Set(fmt.Sprintf("%v_%v", query, category), results, ttl...)
}

func (db DB) Get(k string, o interface{}) error { return nil }

func (db DB) GetResults(query string, category string) ([]result.Result, error) {
	var results []result.Result
	err := db.Get(fmt.Sprintf("%v_%v", query, category), &results)
	return results, err
}

func (db DB) GetTTL(k string) (time.Duration, error) { return 0, nil }

func (db DB) GetResultsTTL(query string, category string) (time.Duration, error) {
	return db.GetTTL(fmt.Sprintf("%v_%v", query, category))
}
