package search

import (
	"sync"

	"github.com/hearchco/agent/src/search/engines"
	"github.com/hearchco/agent/src/search/engines/options"
	"github.com/hearchco/agent/src/search/result"
	"github.com/hearchco/agent/src/search/scraper"
)

// Searchers.
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
			runEngine(groupName, onceWrapMap[engName], concMap, engName, searcher.Search, query, opts)
		}()
	}
}

// Image searchers.
func runRequiredImageSearchers(engs []engines.Name, searchers []scraper.ImageSearcher, wgRequiredEngines *sync.WaitGroup, concMap *result.ResultConcMap, query string, opts options.Options, onceWrapMap map[engines.Name]*onceWrapper) {
	runImageSearchers(groupRequired, engs, searchers, wgRequiredEngines, concMap, query, opts, onceWrapMap)
}

func runPreferredImageSearchers(engs []engines.Name, searchers []scraper.ImageSearcher, wgPreferredEngines *sync.WaitGroup, concMap *result.ResultConcMap, query string, opts options.Options, onceWrapMap map[engines.Name]*onceWrapper) {
	runImageSearchers(groupPreferred, engs, searchers, wgPreferredEngines, concMap, query, opts, onceWrapMap)
}

func runImageSearchers(groupName string, engs []engines.Name, searchers []scraper.ImageSearcher, wgRequiredEngines *sync.WaitGroup, concMap *result.ResultConcMap, query string, opts options.Options, onceWrapMap map[engines.Name]*onceWrapper) {
	wgRequiredEngines.Add(len(engs))
	for _, engName := range engs {
		searcher := searchers[engName]
		go func() {
			// Indicate that the engine is done.
			defer wgRequiredEngines.Done()

			// Run the engine.
			runEngine(groupName, onceWrapMap[engName], concMap, engName, searcher.ImageSearch, query, opts)
		}()
	}
}

// Suggesters.
func runRequiredSuggesters(engs []engines.Name, suggesters []scraper.Suggester, wgRequiredEngines *sync.WaitGroup, concMap *result.SuggestionConcMap, query string, opts options.Options, onceWrapMap map[engines.Name]*onceWrapper) {
	runSuggesters(groupRequired, engs, suggesters, wgRequiredEngines, concMap, query, opts, onceWrapMap)
}

func runPreferredSuggesters(engs []engines.Name, suggesters []scraper.Suggester, wgPreferredEngines *sync.WaitGroup, concMap *result.SuggestionConcMap, query string, opts options.Options, onceWrapMap map[engines.Name]*onceWrapper) {
	runSuggesters(groupPreferred, engs, suggesters, wgPreferredEngines, concMap, query, opts, onceWrapMap)
}

func runSuggesters(groupName string, engs []engines.Name, suggesters []scraper.Suggester, wgRequiredEngines *sync.WaitGroup, concMap *result.SuggestionConcMap, query string, opts options.Options, onceWrapMap map[engines.Name]*onceWrapper) {
	wgRequiredEngines.Add(len(engs))
	for _, engName := range engs {
		suggester := suggesters[engName]
		go func() {
			// Indicate that the engine is done.
			defer wgRequiredEngines.Done()

			// Run the engine.
			runEngine(groupName, onceWrapMap[engName], concMap, engName, suggester.Suggest, query, opts)
		}()
	}
}
