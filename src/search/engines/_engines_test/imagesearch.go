package _engines_test

import (
	"context"
	"strings"
	"testing"

	"github.com/hearchco/agent/src/search/result"
	"github.com/hearchco/agent/src/search/scraper"
)

func CheckImageSearch(t *testing.T, e scraper.ImageSearcher, tchar []TestCaseHasAnyResults, tccr []TestCaseContainsResults, tcrr []TestCaseRankedResults) {
	// TestCaseHasAnyResults
	for _, tc := range tchar {
		e.ReInitSearcher(context.Background())

		resChan := make(chan result.ResultScraped, 100)
		go e.ImageSearch(tc.Query, tc.Options, resChan)

		results := make([]result.ResultScraped, 0)
		for r := range resChan {
			results = append(results, r)
		}

		if len(results) == 0 {
			defer t.Errorf("Got no results for %q", tc.Query)
		}
	}

	// TestCaseContainsResults
	for _, tc := range tccr {
		e.ReInitSearcher(context.Background())

		resChan := make(chan result.ResultScraped, 100)
		go e.ImageSearch(tc.Query, tc.Options, resChan)

		results := make([]result.ResultScraped, 0)
		for r := range resChan {
			results = append(results, r)
		}

		if len(results) == 0 {
			defer t.Errorf("Got no results for %q", tc.Query)
		} else {
			for _, rURL := range tc.ResultURLs {
				found := false

				for _, r := range results {
					if strings.Contains(r.URL(), rURL) {
						found = true
						break
					}
				}

				if !found {
					defer t.Errorf("Couldn't find %q (%q).\nThe results: %q", rURL, tc.Query, results)
				}
			}
		}
	}

	// TestCaseRankedResults
	for _, tc := range tcrr {
		e.ReInitSearcher(context.Background())

		resChan := make(chan result.ResultScraped, 100)
		go e.ImageSearch(tc.Query, tc.Options, resChan)

		results := make([]result.ResultScraped, 0)
		for r := range resChan {
			results = append(results, r)
		}

		if len(results) == 0 {
			defer t.Errorf("Got no results for %q", tc.Query)
		} else if len(results) < len(tc.ResultURLs) {
			defer t.Errorf("Number of results is less than test case URLs.")
		} else {
			for i, rURL := range tc.ResultURLs {
				if !strings.Contains(results[i].URL(), rURL) {
					defer t.Errorf("Wrong result on rank %q: %q (%q).\nThe results: %q", i+1, rURL, tc.Query, results)
				}
			}
		}
	}
}
