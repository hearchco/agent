package duckduckgo_test

import (
	"testing"

	"github.com/tminaorg/brzaguza/src/config"
	"github.com/tminaorg/brzaguza/src/engines"
	"github.com/tminaorg/brzaguza/src/search"
)

type TestCase struct {
	query   string
	options engines.Options
}

func TestSearch(t *testing.T) {
	// testing config
	conf := config.Config{
		Engines: map[string]config.Engine{
			engines.DuckDuckGo.ToLower(): {
				Enabled: true,
			},
		},
	}

	// enabled engines names
	for name, engine := range conf.Engines {
		if engine.Enabled {
			if engineName, err := engines.NameString(name); err == nil {
				config.EnabledEngines = append(config.EnabledEngines, engineName)
			}
		}
	}

	// test cases
	testCases := [...]TestCase{{
		query: "wikipedia",
		options: engines.Options{
			MaxPages:   1,
			VisitPages: false,
		},
	}}

	// running tests
	for _, tc := range testCases {
		if results := search.PerformSearch(tc.query, tc.options, &conf); len(results) == 0 {
			t.Errorf("Got no results for %v", tc.query)
		}
	}
}
