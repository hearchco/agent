package _engines_test

import (
	"strings"
	"testing"

	"github.com/hearchco/agent/src/search/result"
	"github.com/hearchco/agent/src/search/scraper"
)

func CheckSuggest(t *testing.T, e scraper.Suggester, q string) {
	sugChan := make(chan result.SuggestionScraped)
	go func() {
		err, found := e.Suggest(q, NewOpts(), sugChan)
		if len(err) > 0 || !found {
			t.Errorf("Failed to get suggestions: %v", err)
		}
	}()

	suggs := make([]string, 0, 10)
	for sug := range sugChan {
		suggs = append(suggs, sug.Value())
	}
	if len(suggs) == 0 {
		t.Errorf("No suggestions returned")
	}

	for _, s := range suggs {
		if s == "" {
			t.Errorf("Empty suggestion")
		} else if !strings.Contains(s, q) {
			t.Errorf("Suggestion doesn't contain query (%q): %q", q, s)
		}
	}
}
