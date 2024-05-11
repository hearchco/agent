package _engines_test

import (
	"time"

	"github.com/hearchco/hearchco/src/config"
	"github.com/hearchco/hearchco/src/search/category"
	"github.com/hearchco/hearchco/src/search/engines"
)

type TestCaseHasAnyResults struct {
	Query   string
	Options engines.Options
}

type TestCaseContainsResults struct {
	Query     string
	ResultURL []string
	Options   engines.Options
}

type TestCaseRankedResults struct {
	Query     string
	ResultURL []string
	Options   engines.Options
}

func NewConfig(engineName engines.Name) config.Config {
	return config.Config{
		Categories: map[category.Name]config.Category{
			category.GENERAL: {
				Engines: []engines.Name{engineName},
				Ranking: config.NewRanking(),
				Timings: config.CategoryTimings{
					HardTimeout: 10000 * time.Millisecond,
				},
			},
			category.IMAGES: {
				Engines: []engines.Name{engineName},
				Ranking: config.NewRanking(),
				Timings: config.CategoryTimings{
					HardTimeout: 10000 * time.Millisecond,
				},
			},
		},
	}
}

func NewOpts() engines.Options {
	return engines.Options{
		Pages: engines.Pages{Start: 0, Max: 1},
	}
}
