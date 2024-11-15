package search

import (
	"context"

	"github.com/hearchco/agent/src/search/engines"
	"github.com/hearchco/agent/src/search/scraper"
)

// Initialize web searchers.
func initializeWebSearchers(ctx context.Context, engs []engines.Name) []scraper.WebSearcher {
	searchers := webSearcherArray()
	for _, engName := range engs {
		searchers[engName].InitSearcher(ctx)
	}
	return searchers[:]
}

// Initialize image searchers.
func initializeImageSearchers(ctx context.Context, engs []engines.Name) []scraper.ImageSearcher {
	searchers := imageSearcherArray()
	for _, engName := range engs {
		searchers[engName].InitSearcher(ctx)
	}
	return searchers[:]
}

// Initialize suggesters.
func initializeSuggesters(ctx context.Context, engs []engines.Name) []scraper.Suggester {
	suggesters := suggesterArray()
	for _, engName := range engs {
		suggesters[engName].InitSuggester(ctx)
	}
	return suggesters[:]
}
