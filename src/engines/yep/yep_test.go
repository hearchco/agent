package yep_test

import (
	"testing"

	"github.com/tminaorg/brzaguza/src/engines"
	"github.com/tminaorg/brzaguza/src/engines/_engines_test"
)

func TestSearch(t *testing.T) {
	engineName := engines.YEP

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
		Query:     "youtube",
		ResultURL: []string{"youtube.com"},
		Options: engines.Options{
			MaxPages:   1,
			VisitPages: false,
		},
	}}

	tcrr := [...]_engines_test.TestCaseRankedResults{{
		Query:     "wikipedia",
		ResultURL: []string{"wikipedia."},
		Options: engines.Options{
			MaxPages:   1,
			VisitPages: false,
		},
	}}

	_engines_test.CheckTestCases(tchar[:], tccr[:], tcrr[:], t, conf)
}
