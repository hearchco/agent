package _engines_test

import (
	"github.com/tminaorg/brzaguza/src/category"
	"github.com/tminaorg/brzaguza/src/config"
	"github.com/tminaorg/brzaguza/src/engines"
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
			},
		},
	}
}
