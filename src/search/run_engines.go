package search

import (
	"sync"

	"github.com/hearchco/agent/src/search/engines"
	"github.com/hearchco/agent/src/search/engines/options"
	"github.com/hearchco/agent/src/search/result"
	"github.com/hearchco/agent/src/search/scraper"
)

func runRequiredSearchers(engs []engines.Name, searchers []scraper.Searcher, wgRequiredEngines *sync.WaitGroup, concMap *result.ResultConcMap, query string, opts options.Options, onceWrapMap map[engines.Name]*onceWrapper) {
	runSearchers(groupRequired, engs, searchers, wgRequiredEngines, concMap, query, opts, onceWrapMap)
}

func runPreferredSearchers(engs []engines.Name, searchers []scraper.Searcher, wgPreferredEngines *sync.WaitGroup, concMap *result.ResultConcMap, query string, opts options.Options, onceWrapMap map[engines.Name]*onceWrapper) {
	runSearchers(groupPreferred, engs, searchers, wgPreferredEngines, concMap, query, opts, onceWrapMap)
}

func runSearchers(groupName string, engs []engines.Name, searchers []scraper.Searcher, wgRequiredEngines *sync.WaitGroup, concMap *result.ResultConcMap, query string, opts options.Options, onceWrapMap map[engines.Name]*onceWrapper) {
	wgRequiredEngines.Add(len(engs))
	for _, engName := range engs {
		searcher := searchers[engName]
		go func() {
			// Indicate that the engine is done.
			defer wgRequiredEngines.Done()

			// Run the engine.
			runSearcher(groupName, onceWrapMap[engName], concMap, engName, searcher, query, opts)
		}()
	}
}

func runRequiredSuggesters(engs []engines.Name, suggesters []scraper.Suggester, wgRequiredEngines *sync.WaitGroup, concMap *result.SuggestionConcMap, query string, locale options.Locale, onceWrapMap map[engines.Name]*onceWrapper) {
	runSuggesters(groupRequired, engs, suggesters, wgRequiredEngines, concMap, query, locale, onceWrapMap)
}

func runPreferredSuggesters(engs []engines.Name, suggesters []scraper.Suggester, wgPreferredEngines *sync.WaitGroup, concMap *result.SuggestionConcMap, query string, locale options.Locale, onceWrapMap map[engines.Name]*onceWrapper) {
	runSuggesters(groupPreferred, engs, suggesters, wgPreferredEngines, concMap, query, locale, onceWrapMap)
}

func runSuggesters(groupName string, engs []engines.Name, suggesters []scraper.Suggester, wgRequiredEngines *sync.WaitGroup, concMap *result.SuggestionConcMap, query string, locale options.Locale, onceWrapMap map[engines.Name]*onceWrapper) {
	wgRequiredEngines.Add(len(engs))
	for _, engName := range engs {
		searcher := suggesters[engName]
		go func() {
			// Indicate that the engine is done.
			defer wgRequiredEngines.Done()

			// Run the engine.
			runSuggester(groupName, onceWrapMap[engName], concMap, engName, searcher, query, locale)
		}()
	}
}
