package cache

import (
	"time"

	"github.com/hearchco/hearchco/src/search/engines"
	"github.com/hearchco/hearchco/src/search/result"
)

func (db DB) SetResults(query string, options engines.Options, results []result.Result, ttl ...time.Duration) error {
	key := combineQueryWithOptions(query, options)
	return db.driver.Set(key, results, ttl...)
}

func (db DB) GetResults(query string, options engines.Options) ([]result.Result, error) {
	key := combineQueryWithOptions(query, options)
	var results []result.Result
	err := db.driver.Get(key, &results)
	return results, err
}

func (db DB) GetResultsTTL(query string, options engines.Options) (time.Duration, error) {
	key := combineQueryWithOptions(query, options)
	return db.driver.GetTTL(key)
}

func combineQueryWithOptions(query string, options engines.Options) string {
	return combineIntoKey(query, options.VisitPages, options.SafeSearch, options.Pages.Start, options.Pages.Max, options.Locale, options.Category.String())
}
