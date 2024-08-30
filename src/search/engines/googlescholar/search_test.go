package googlescholar

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
		Query:      "interaction nets",
		ResultURLs: []string{"https://dl.acm.org/doi/pdf/10.1145/96709.96718"},
		Options:    opt,
	}}

	tcrr := []_engines_test.TestCaseRankedResults{{
		Query:      "On building fast kd-trees for ray tracing, and on doing that in O (N log N)",
		ResultURLs: []string{"https://ieeexplore.ieee.org/abstract/document/4061547/"},
		Options:    opt,
	}}

	se := New()
	se.InitSearcher(context.Background(), conf.Categories[category.GENERAL].Timings)

	_engines_test.CheckSearch(t, se, tchar, tccr, tcrr)
}
