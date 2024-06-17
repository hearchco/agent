package _engines_test

import (
	"time"

	"github.com/hearchco/agent/src/config"
	"github.com/hearchco/agent/src/search/category"
	"github.com/hearchco/agent/src/search/engines"
	"github.com/hearchco/agent/src/search/engines/options"
)

type TestCaseHasAnyResults struct {
	Query   string
	Options options.Options
}

type TestCaseContainsResults struct {
	Query      string
	ResultURLs []string
	Options    options.Options
}

type TestCaseRankedResults struct {
	Query      string
	ResultURLs []string
	Options    options.Options
}

func NewConfig(seName engines.Name) config.Config {
	return config.Config{
		Categories: map[category.Name]config.Category{
			category.GENERAL: {
				Engines: []engines.Name{seName},
				Ranking: config.EmptyRanking([]engines.Name{seName}),
				Timings: config.CategoryTimings{
					HardTimeout: 10000 * time.Millisecond,
				},
			},
			category.IMAGES: {
				Engines: []engines.Name{seName},
				Ranking: config.EmptyRanking([]engines.Name{seName}),
				Timings: config.CategoryTimings{
					HardTimeout: 10000 * time.Millisecond,
				},
			},
		},
	}
}

func NewOpts() options.Options {
	return options.Options{
		Pages:      options.Pages{Start: 0, Max: 1},
		Locale:     options.LocaleDefault,
		SafeSearch: false,
	}
}
