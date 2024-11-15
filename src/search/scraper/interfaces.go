package scraper

import (
	"context"

	"github.com/hearchco/agent/src/search/engines"
	"github.com/hearchco/agent/src/search/engines/options"
	"github.com/hearchco/agent/src/search/result"
)

// Interface that each search engine must implement to be a Search Engine.
type Enginer interface {
	GetName() engines.Name
	GetOrigins() []engines.Name
	Init(context.Context)
}

// Interface that each search engine must implement to support searching web results.
type WebSearcher interface {
	Enginer

	InitSearcher(context.Context)
	WebSearch(string, options.Options, chan result.ResultScraped) ([]error, bool)
}

// Interface that each search engine must implement to support searching image results.
type ImageSearcher interface {
	Enginer

	InitSearcher(context.Context)
	ImageSearch(string, options.Options, chan result.ResultScraped) ([]error, bool)
}

// Interface that each search engine must implement to support suggesting.
type Suggester interface {
	Enginer

	InitSuggester(context.Context)
	Suggest(string, options.Options, chan result.SuggestionScraped) ([]error, bool)
}
