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
}

// Interface that each search engine must implement to support searching general results.
type Searcher interface {
	Enginer

	InitSearcher(context.Context, config.CategoryTimings)
	ReInitSearcher(context.Context)
	Search(string, options.Options, chan result.ResultScraped) ([]error, bool)
}

// Interface that each search engine must implement to support searching image results.
type ImageSearcher interface {
	Enginer

	InitSearcher(context.Context, config.CategoryTimings)
	ReInitSearcher(context.Context)
	ImageSearch(string, options.Options, chan result.ResultScraped) ([]error, bool)
}

// Interface that each search engine must implement to support suggesting.
type Suggester interface {
	Enginer

	InitSuggester(context.Context, config.CategoryTimings)
	ReInitSuggester(context.Context)
	Suggest(string, options.Options, chan result.SuggestionScraped) ([]error, bool)
}
