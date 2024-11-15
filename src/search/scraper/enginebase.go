package scraper

import (
	"context"

	"github.com/gocolly/colly/v2"

	"github.com/hearchco/agent/src/search/engines"
)

// Base struct for every search engine.
type EngineBase struct {
	Name      engines.Name
	Origins   []engines.Name
	collector *colly.Collector
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
func (e *EngineBase) Init(ctx context.Context) {
	e.initCollectorOnRequest(ctx)
	e.initCollectorOnResponse()
	e.initCollectorOnError()
}

// Used to initialize the EngineBase collector for searching.
func (e *EngineBase) InitSearcher(ctx context.Context) {
	e.initCollectorSearcher(ctx)
	e.Init(ctx)
}

// Used to initialize the EngineBase collector for suggesting.
func (e *EngineBase) InitSuggester(ctx context.Context) {
	e.initCollectorSuggester(ctx)
	e.Init(ctx)
}
