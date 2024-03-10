package googleimages_test

import (
	"testing"

	"github.com/hearchco/hearchco/src/search/engines"
	"github.com/hearchco/hearchco/src/search/engines/_engines_test"
)

func TestSearch(t *testing.T) {
	engineName := engines.GOOGLEIMAGES

	// testing config
	conf := _engines_test.NewConfig(engineName)
	opt := _engines_test.NewOpts()

	// test cases
	tchar := [...]_engines_test.TestCaseHasAnyResults{{
		Query:   "ping",
		Options: opt,
	}}

	tccr := [...]_engines_test.TestCaseContainsResults{{
		Query:     "wikipedia logo",
		ResultURL: []string{"upload.wikimedia.org"},
		Options:   opt,
	}}

	tcrr := [...]_engines_test.TestCaseRankedResults{{
		Query:     "linux logo wikipedia",
		ResultURL: []string{"upload.wikimedia.org"},
		Options:   opt,
	}}

	_engines_test.CheckTestCases(tchar[:], tccr[:], tcrr[:], t, conf)
}
