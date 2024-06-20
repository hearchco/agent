package scraper

import (
	"context"

	"github.com/gocolly/colly/v2"

	"github.com/hearchco/agent/src/config"
	"github.com/hearchco/agent/src/search/engines"
)

// Base struct for every search engine.
type EngineBase struct {
	Name      engines.Name
	Origins   []engines.Name
	collector *colly.Collector
	timings   config.CategoryTimings
}

// Used to get the name of the search engine.
func (e EngineBase) GetName() engines.Name {
	return e.Name
}

// Used to get the origins of the search engine.
func (e EngineBase) GetOrigins() []engines.Name {
	return e.Origins
}

// Used to initialize the EngineBase collector.
func (e *EngineBase) Init(ctx context.Context, timings config.CategoryTimings) {
	e.timings = timings
	e.initCollector(ctx)
	e.initLimitRule(timings)
	e.initCollectorOnRequest(ctx)
	e.initCollectorOnResponse()
	e.initCollectorOnError()
}

// Used to allow re-running the Search method.
func (e *EngineBase) ReInit(ctx context.Context) {
	e.Init(ctx, e.timings)
}
