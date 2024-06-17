package cache

import (
	"fmt"
	"time"

	"github.com/hearchco/agent/src/search/category"
	"github.com/hearchco/agent/src/search/engines/options"
	"github.com/hearchco/agent/src/search/result"
)

func (db DB) SetResults(q string, cat category.Name, opts options.Options, results []result.Result, ttl ...time.Duration) error {
	key := combineQueryWithOptions(q, cat, opts)
	return db.driver.Set(key, results, ttl...)
}

func (db DB) GetResults(q string, cat category.Name, opts options.Options) ([]result.Result, error) {
	var results []result.Result
	var err error

	key := combineQueryWithOptions(q, cat, opts)
	if cat == category.IMAGES {
		var imgResults []result.Images
		err = db.driver.Get(key, &imgResults)
		results = make([]result.Result, 0, len(imgResults))
		for _, imgResult := range imgResults {
			results = append(results, &imgResult)
		}

	} else {
		var genResults []result.General
		err = db.driver.Get(key, &genResults)
		results = make([]result.Result, 0, len(genResults))
		for _, imgResult := range genResults {
			results = append(results, &imgResult)
		}
	}

	return results, err
}

func (db DB) GetResultsTTL(q string, cat category.Name, opts options.Options) (time.Duration, error) {
	key := combineQueryWithOptions(q, cat, opts)
	return db.driver.GetTTL(key)
}

func combineQueryWithOptions(q string, cat category.Name, opts options.Options) string {
	return combineIntoKey(q, cat.String(), opts.Pages.Start, opts.Pages.Max, opts.Locale, opts.SafeSearch)
}

func combineIntoKey(s ...any) string {
	var key string
	for i, v := range s {
		if i == 0 {
			key = fmt.Sprintf("%v", v)
		} else {
			key = fmt.Sprintf("%v_%v", key, v)
		}
	}
	return key
}
