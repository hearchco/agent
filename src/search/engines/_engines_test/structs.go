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

func NewConfig(engineName engines.Name) *config.Config {
	config.EnabledEngines = append(config.EnabledEngines, engineName)
	return &config.Config{
		Categories: map[category.Name]config.Category{
			category.GENERAL: {
				Engines: []engines.Name{engineName},
				Ranking: config.NewRanking(),
				Timings: config.Timings{
					Timeout: 10000 * time.Millisecond, // colly default
				},
			},
		},
	}
}
