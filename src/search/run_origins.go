package search

import (
	"slices"
	"sync"

	"github.com/hearchco/agent/src/search/engines"
	"github.com/hearchco/agent/src/search/engines/options"
	"github.com/hearchco/agent/src/search/result"
	"github.com/hearchco/agent/src/search/scraper"
)

// Run required engines.
func runRequiredByOriginSearchers(engs []engines.Name, searchers []scraper.Searcher, wgByOriginEngines *sync.WaitGroup, concMap *result.ResultConcMap, enabledEngines []engines.Name, query string, opts options.Options, onceWrapMap map[engines.Name]*onceWrapper) {
	runByOriginSearchers(groupRequiredByOrigin, engs, searchers, wgByOriginEngines, concMap, enabledEngines, query, opts, onceWrapMap)
}

// Run preferred engines.
func runPreferredByOriginSearchers(engs []engines.Name, searchers []scraper.Searcher, wgByOriginEngines *sync.WaitGroup, concMap *result.ResultConcMap, enabledEngines []engines.Name, query string, opts options.Options, onceWrapMap map[engines.Name]*onceWrapper) {
	runByOriginSearchers(groupPreferredByOrigin, engs, searchers, wgByOriginEngines, concMap, enabledEngines, query, opts, onceWrapMap)
}

// Run passed engines by origin.
func runByOriginSearchers(groupName string, engs []engines.Name, searchers []scraper.Searcher, wg *sync.WaitGroup, concMap *result.ResultConcMap, enabledEngines []engines.Name, query string, opts options.Options, onceWrapMap map[engines.Name]*onceWrapper) {
	// Create a map of slices of all the engines that contain origins from the engines by origin.
	engsMap := make(map[engines.Name][]engines.Name, len(engs))
	for _, originName := range engs {
		for _, engName := range enabledEngines {
			origins := searchers[engName].GetOrigins()
			if slices.Contains(origins, originName) {
				workers, ok := engsMap[originName]
				if !ok {
					workers = make([]engines.Name, 0, len(enabledEngines))
				}
				engsMap[originName] = append(workers, engName)
			}
		}
	}

	// Run all by origin engines. Cond should be awaited unless the timeout is reached.
	wg.Add(len(engsMap))
	for _, workers := range engsMap {
		if len(workers) == 0 {
			wg.Done()
			continue
		}

		var wgWorkers sync.WaitGroup
		wgWorkers.Add(len(workers))
		successOrigin := sync.Cond{L: &sync.Mutex{}}
		go waitForSuccessOrFinish(&successOrigin, &wgWorkers, wg)

		for _, engName := range workers {
			searcher := searchers[engName]
			go func() {
				// Indicate that the engine is done (successful or not).
				defer wgWorkers.Done()

				// Run the engine.
				runSearcher(groupName, onceWrapMap[engName], concMap, engName, searcher, query, opts)

				// Indicate that the engine was successful.
				if onceWrapMap[engName].Success() {
					successOrigin.L.Lock()
					successOrigin.Signal()
					successOrigin.L.Unlock()
				}
			}()
		}
	}
}

// Run required engines.
func runRequiredByOriginSuggesters(engs []engines.Name, suggesters []scraper.Suggester, wgByOriginEngines *sync.WaitGroup, concMap *result.SuggestionConcMap, enabledEngines []engines.Name, query string, locale options.Locale, onceWrapMap map[engines.Name]*onceWrapper) {
	runByOriginSuggesters(groupRequiredByOrigin, engs, suggesters, wgByOriginEngines, concMap, enabledEngines, query, locale, onceWrapMap)
}

// Run preferred engines.
func runPreferredByOriginSuggesters(engs []engines.Name, suggesters []scraper.Suggester, wgByOriginEngines *sync.WaitGroup, concMap *result.SuggestionConcMap, enabledEngines []engines.Name, query string, locale options.Locale, onceWrapMap map[engines.Name]*onceWrapper) {
	runByOriginSuggesters(groupPreferredByOrigin, engs, suggesters, wgByOriginEngines, concMap, enabledEngines, query, locale, onceWrapMap)
}

// Run passed engines by origin.
func runByOriginSuggesters(groupName string, engs []engines.Name, suggesters []scraper.Suggester, wg *sync.WaitGroup, concMap *result.SuggestionConcMap, enabledEngines []engines.Name, query string, locale options.Locale, onceWrapMap map[engines.Name]*onceWrapper) {
	// Create a map of slices of all the engines that contain origins from the engines by origin.
	engsMap := make(map[engines.Name][]engines.Name, len(engs))
	for _, originName := range engs {
		for _, engName := range enabledEngines {
			origins := suggesters[engName].GetOrigins()
			if slices.Contains(origins, originName) {
				workers, ok := engsMap[originName]
				if !ok {
					workers = make([]engines.Name, 0, len(enabledEngines))
				}
				engsMap[originName] = append(workers, engName)
			}
		}
	}

	// Run all by origin engines. Cond should be awaited unless the timeout is reached.
	wg.Add(len(engsMap))
	for _, workers := range engsMap {
		if len(workers) == 0 {
			wg.Done()
			continue
		}

		var wgWorkers sync.WaitGroup
		wgWorkers.Add(len(workers))
		successOrigin := sync.Cond{L: &sync.Mutex{}}
		go waitForSuccessOrFinish(&successOrigin, &wgWorkers, wg)

		for _, engName := range workers {
			suggester := suggesters[engName]
			go func() {
				// Indicate that the engine is done (successful or not).
				defer wgWorkers.Done()

				// Run the engine.
				runSuggester(groupName, onceWrapMap[engName], concMap, engName, suggester, query, locale)

				// Indicate that the engine was successful.
				if onceWrapMap[engName].Success() {
					successOrigin.L.Lock()
					successOrigin.Signal()
					successOrigin.L.Unlock()
				}
			}()
		}
	}
}

// Waits on either c.Wait() or wg.Wait() to do final.Done().
func waitForSuccessOrFinish(c *sync.Cond, wg *sync.WaitGroup, final *sync.WaitGroup) {
	defer final.Done()
	d := sync.Cond{L: &sync.Mutex{}}

	// Wait for signal from any successful worker.
	go func() {
		c.L.Lock()
		c.Wait()
		c.L.Unlock()

		d.L.Lock()
		d.Signal()
		d.L.Unlock()
	}()

	// Wait for all workers to finish (even if it's unsuccessful).
	go func() {
		wg.Wait()

		d.L.Lock()
		d.Signal()
		d.L.Unlock()
	}()

	// Whichever of the above finishes first, signal the final wait group.
	d.L.Lock()
	d.Wait()
	d.L.Unlock()
}
