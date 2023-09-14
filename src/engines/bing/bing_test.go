package bing_test

import (
	"testing"

	"github.com/tminaorg/brzaguza/src/engines"
	"github.com/tminaorg/brzaguza/src/engines/engines_test"
)

func TestSearch(t *testing.T) {
	engineName := engines.Bing

	// testing config
	conf := engines_test.NewConfig(engineName)

	// test cases
	tchar := [...]engines_test.TestCaseHasAnyResults{{
		Query: "ping",
		Options: engines.Options{
			MaxPages:   1,
			VisitPages: false,
		},
	}}

	tccr := [...]engines_test.TestCaseContainsResults{{
		Query:     "facebook",
		ResultURL: []string{"facebook.com"},
		Options: engines.Options{
			MaxPages:   1,
			VisitPages: false,
		},
	}}

	tcrr := [...]engines_test.TestCaseRankedResults{{
		Query:     "wikipedia",
		ResultURL: []string{"wikipedia."},
		Options: engines.Options{
			MaxPages:   1,
			VisitPages: false,
		},
	}}

	engines_test.CheckTestCases(tchar[:], tccr[:], tcrr[:], t, conf)
}
