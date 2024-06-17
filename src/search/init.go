package search

import (
	"context"

	"github.com/hearchco/agent/src/config"
	"github.com/hearchco/agent/src/search/engines"
	"github.com/hearchco/agent/src/search/scraper"
)

func initializeEnginers(ctx context.Context, engs []engines.Name, timings config.CategoryTimings) []scraper.Enginer {
	enginers := enginerArray()
	for _, engName := range engs {
		enginers[engName].Init(ctx, timings)
	}
	return enginers[:]
}
