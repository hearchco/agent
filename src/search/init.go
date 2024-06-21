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
		searchers[engName].Init(ctx, timings)
	}
	return searchers[:]
}

func initializeSuggesters(ctx context.Context, timings config.CategoryTimings) []scraper.Suggester {
	suggesters := suggesterArray()
	for _, suggester := range suggesters {
		suggester.Init(ctx, timings)
	}
	return suggesters[:]
}
