package search

import (
	"context"

	"github.com/hearchco/agent/src/search/engines"
	"github.com/hearchco/agent/src/search/scraper"
)

// Searchers.
func initializeSearchers(ctx context.Context, engs []engines.Name) []scraper.Searcher {
	searchers := searcherArray()
	for _, engName := range engs {
		searchers[engName].InitSearcher(ctx)
	}
	return searchers[:]
}

// Image searchers.
func initializeImageSearchers(ctx context.Context, engs []engines.Name) []scraper.ImageSearcher {
	searchers := imageSearcherArray()
	for _, engName := range engs {
		searchers[engName].InitSearcher(ctx)
	}
	return searchers[:]
}

// Suggesters.
func initializeSuggesters(ctx context.Context, engs []engines.Name) []scraper.Suggester {
	suggesters := suggesterArray()
	for _, engName := range engs {
		suggesters[engName].InitSuggester(ctx)
	}
	return suggesters[:]
}
