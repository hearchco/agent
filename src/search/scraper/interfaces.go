package scraper

import (
	"context"

	"github.com/hearchco/agent/src/config"
	"github.com/hearchco/agent/src/search/engines"
	"github.com/hearchco/agent/src/search/engines/options"
	"github.com/hearchco/agent/src/search/result"
)

// Interface that each search engine must implement to be a Search Engine.
type Enginer interface {
	GetName() engines.Name
	GetOrigins() []engines.Name
	Init(context.Context, config.CategoryTimings)
	ReInit(context.Context)
}

// Interface that each search engine must implement to support searching.
type Searcher interface {
	Enginer

	Search(string, options.Options, chan result.ResultScraped) ([]error, bool)
}

// Interface that each search engine must implement to support suggesting.
type Suggester interface {
	Enginer

	InitSuggest(ctx context.Context, timings config.CategoryTimings)
	Suggest(string, options.Locale, chan []result.SuggestionScraped) (error, bool)
}
