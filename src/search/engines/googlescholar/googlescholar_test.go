package googlescholar_test

import (
	"testing"

	"github.com/hearchco/hearchco/src/search/engines"
	"github.com/hearchco/hearchco/src/search/engines/_engines_test"
)

func TestSearch(t *testing.T) {
	engineName := engines.GOOGLESCHOLAR

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
		Query:     "interaction nets",
		ResultURL: []string{"https://dl.acm.org/doi/pdf/10.1145/96709.96718"},
		Options: engines.Options{
			MaxPages:   1,
			VisitPages: false,
		},
	}}

	tcrr := [...]_engines_test.TestCaseRankedResults{{
		Query:     "On building fast kd-trees for ray tracing, and on doing that in O (N log N)",
		ResultURL: []string{"https://ieeexplore.ieee.org/abstract/document/4061547/"},
		Options: engines.Options{
			MaxPages:   1,
			VisitPages: false,
		},
	}}

	_engines_test.CheckTestCases(tchar[:], tccr[:], tcrr[:], t, conf)
}
