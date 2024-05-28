package cache

import (
	"time"

	"github.com/hearchco/hearchco/src/search/category"
	"github.com/hearchco/hearchco/src/search/result"
)

func (db DB) SetResults(query string, category category.Name, results []result.Result, ttl ...time.Duration) error {
	key := combineIntoKey(query, category.String())
	return db.driver.Set(key, results, ttl...)
}

func (db DB) GetResults(query string, category category.Name) ([]result.Result, error) {
	key := combineIntoKey(query, category.String())
	var results []result.Result
	err := db.driver.Get(key, &results)
	return results, err
}

func (db DB) GetResultsTTL(query string, category category.Name) (time.Duration, error) {
	key := combineIntoKey(query, category.String())
	return db.driver.GetTTL(key)
}
