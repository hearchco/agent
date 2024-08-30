package search

import (
	"context"

	"github.com/hearchco/agent/src/config"
	"github.com/hearchco/agent/src/search/engines"
	"github.com/hearchco/agent/src/search/scraper"
)

// Searchers.
func initializeSearchers(ctx context.Context, engs []engines.Name, timings config.CategoryTimings) []scraper.Searcher {
	searchers := searcherArray()
	for _, engName := range engs {
		searchers[engName].InitSearcher(ctx, timings)
	}
	return searchers[:]
}

// Image searchers.
func initializeImageSearchers(ctx context.Context, engs []engines.Name, timings config.CategoryTimings) []scraper.ImageSearcher {
	searchers := imageSearcherArray()
	for _, engName := range engs {
		searchers[engName].InitSearcher(ctx, timings)
	}
	return searchers[:]
}

// Suggesters.
func initializeSuggesters(ctx context.Context, engs []engines.Name, timings config.CategoryTimings) []scraper.Suggester {
	suggesters := suggesterArray()
	for _, engName := range engs {
		suggesters[engName].InitSuggester(ctx, timings)
	}
	return suggesters[:]
}
