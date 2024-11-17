package google

import (
	"context"
	"testing"

	"github.com/hearchco/agent/src/search/engines/_engines_test"
)

func TestImageSearch(t *testing.T) {
	// Testing options.
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
		ResultURLs: []string{"upload.wikimedia.org"},
		Options:    opt,
	}}

	se := New()
	se.InitSearcher(context.Background())

	_engines_test.CheckImageSearch(t, se, tchar[:], tccr[:], tcrr[:])
}
