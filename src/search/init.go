package search

import (
	"context"

	"github.com/hearchco/agent/src/config"
	"github.com/hearchco/agent/src/search/engines"
	"github.com/hearchco/agent/src/search/scraper"
)

func initializeSearchers(ctx context.Context, engs []engines.Name, timings config.CategoryTimings) []scraper.Searcher {
	searchers := searcherArray()
	for _, engName := range engs {
		searchers[engName].InitSearcher(ctx, timings)
	}
	return searchers[:]
}

func initializeSuggesters(ctx context.Context, engs []engines.Name, timings config.CategoryTimings) []scraper.Suggester {
	suggesters := suggesterArray()
	for _, engName := range engs {
		suggesters[engName].InitSuggester(ctx, timings)
	}
	return suggesters[:]
}
