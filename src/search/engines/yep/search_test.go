package yep

// import (
// 	"context"
// 	"testing"

// 	"github.com/hearchco/agent/src/search/category"
// 	"github.com/hearchco/agent/src/search/engines/_engines_test"
// )

// func TestSearch(t *testing.T) {
// 	// Testing options.
// 	conf := _engines_test.NewConfig(seName)
// 	opt := _engines_test.NewOpts()

// 	// Test cases.
// 	tchar := []_engines_test.TestCaseHasAnyResults{{
// 		Query:   "ping",
// 		Options: opt,
// 	}}

// 	tccr := []_engines_test.TestCaseContainsResults{{
// 		Query:      "youtube",
// 		ResultURLs: []string{"youtube.com"},
// 		Options:    opt,
// 	}}

// 	tcrr := []_engines_test.TestCaseRankedResults{{
// 		Query:      "wikipedia",
// 		ResultURLs: []string{"wikipedia."},
// 		Options:    opt,
// 	}}

// 	se := New()
// 	se.Init(context.Background(), conf.Categories[category.GENERAL].Timings)

// 	_engines_test.CheckSearch(t, se, tchar, tccr, tcrr)
// }
