package bingimages_test

import (
	"testing"

	"github.com/hearchco/hearchco/src/search/engines"
	"github.com/hearchco/hearchco/src/search/engines/_engines_test"
)

func TestSearch(t *testing.T) {
	engineName := engines.BINGIMAGES

	// testing config
	conf := _engines_test.NewConfig(engineName)

	// test cases
	tchar := [...]_engines_test.TestCaseHasAnyResults{{
		Query: "ping",
		Options: engines.Options{
			MaxPages:   1,
			VisitPages: false,
		},
	}}

	tccr := [...]_engines_test.TestCaseContainsResults{{
		Query:     "wikipedia logo",
		ResultURL: []string{"upload.wikimedia.org"},
		Options: engines.Options{
			MaxPages:   1,
			VisitPages: false,
		},
	}}

	tcrr := [...]_engines_test.TestCaseRankedResults{{
		Query:     "linux logo wikipedia",
		ResultURL: []string{"logos-world.net"},
		Options: engines.Options{
			MaxPages:   1,
			VisitPages: false,
		},
	}}

	_engines_test.CheckTestCases(tchar[:], tccr[:], tcrr[:], t, conf)
}
