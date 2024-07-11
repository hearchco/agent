package bingimages

import (
	"context"
	"testing"

	"github.com/hearchco/agent/src/search/category"
	"github.com/hearchco/agent/src/search/engines/_engines_test"
)

func TestSearch(t *testing.T) {
	// Testing options.
	conf := _engines_test.NewConfig(seName)
	opt := _engines_test.NewOpts()

	// Test cases.
	tchar := []_engines_test.TestCaseHasAnyResults{{
		Query:   "ping",
		Options: opt,
	}}

	tccr := []_engines_test.TestCaseContainsResults{{
		Query:      "wikipedia logo",
		ResultURLs: []string{"upload.wikimedia.org"},
		Options:    opt,
	}}

	tcrr := []_engines_test.TestCaseRankedResults{{
		Query:      "linux logo wikipedia",
		ResultURLs: []string{"logos-world.net"},
		Options:    opt,
	}}

	se := New()
	se.InitSearcher(context.Background(), conf.Categories[category.GENERAL].Timings)

	_engines_test.CheckTestCases(t, se, tchar[:], tccr[:], tcrr[:])
}
