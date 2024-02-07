package _engines_test

import (
	"strings"
	"testing"

	"github.com/hearchco/hearchco/src/config"
	"github.com/hearchco/hearchco/src/search"
)

func CheckTestCases(tchar []TestCaseHasAnyResults, tccr []TestCaseContainsResults,
	tcrr []TestCaseRankedResults, t *testing.T, conf *config.Config) {

	// TestCaseHasAnyResults
	for _, tc := range tchar {
		if results := search.PerformSearch(tc.Query, tc.Options, conf); len(results) == 0 {
			defer t.Errorf("Got no results for %v", tc.Query)
		}
	}

	// TestCaseContainsResults
	for _, tc := range tccr {
		results := search.PerformSearch(tc.Query, tc.Options, conf)
		if len(results) == 0 {
			defer t.Errorf("Got no results for %v", tc.Query)
		} else {
			for _, rURL := range tc.ResultURL {
				found := false

				for _, r := range results {
					if strings.Contains(r.URL, rURL) {
						found = true
						break
					}
				}

				if !found {
					defer t.Errorf("Couldn't find %v (%v).\nThe results: %v", rURL, tc.Query, results)
				}
			}
		}
	}

	// TestCaseRankedResults
	for _, tc := range tcrr {
		results := search.PerformSearch(tc.Query, tc.Options, conf)
		if len(results) == 0 {
			defer t.Errorf("Got no results for %v", tc.Query)
		} else if len(results) < len(tc.ResultURL) {
			defer t.Errorf("Number of results is less than test case URLs.")
		} else {
			for i, rURL := range tc.ResultURL {
				if !strings.Contains(results[i].URL, rURL) {
					defer t.Errorf("Wrong result on rank %v: %v (%v).\nThe results: %v", i+1, rURL, tc.Query, results)
				}
			}
		}
	}
}
