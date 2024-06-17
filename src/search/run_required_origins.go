package search

import (
	"slices"
	"sync"

	"github.com/rs/zerolog/log"

	"github.com/hearchco/agent/src/search/engines"
	"github.com/hearchco/agent/src/search/engines/options"
	"github.com/hearchco/agent/src/search/result"
	"github.com/hearchco/agent/src/search/scraper"
	"github.com/hearchco/agent/src/utils/anonymize"
)

func runRequiredByOriginEngines(enginers []scraper.Enginer, wgRequiredByOriginEngines *sync.WaitGroup, query string, opts options.Options, requiredByOriginEngines []engines.Name, enabledEngines []engines.Name, engChan chan chan result.ResultScraped, searchOnce map[engines.Name]*onceWrapper) {
	// Create a map of slices of all the engines that contain origins from the required engines by origin.
	requiredByOriginEnginesMap := make(map[engines.Name][]engines.Name, len(requiredByOriginEngines))
	for _, originName := range requiredByOriginEngines {
		for _, engName := range enabledEngines {
			origins := enginers[engName].GetOrigins()
			if slices.Contains(origins, originName) {
				workers, ok := requiredByOriginEnginesMap[originName]
				if !ok {
					workers = make([]engines.Name, 0, len(enabledEngines))
				}
				requiredByOriginEnginesMap[originName] = append(workers, engName)
			}
		}
	}

	// Run all required by origin engines. Cond should be awaited unless the hard timeout is reached.
	wgRequiredByOriginEngines.Add(len(requiredByOriginEnginesMap))
	for _, workers := range requiredByOriginEnginesMap {
		if len(workers) == 0 {
			wgRequiredByOriginEngines.Done()
			continue
		}

		var wgWorkers sync.WaitGroup
		wgWorkers.Add(len(workers))
		successOrigin := sync.Cond{L: &sync.Mutex{}}
		go waitForSuccessOrFinish(&successOrigin, &wgWorkers, wgRequiredByOriginEngines)

		for _, engName := range workers {
			enginer := enginers[engName]
			resChan := make(chan result.ResultScraped, 100)
			engChan <- resChan
			go func() {
				// Indicate that the worker is done, successful or not.
				defer wgWorkers.Done()

				searchOnce[engName].Do(func() {
					log.Trace().
						Str("engine", engName.String()).
						Str("query", anonymize.String(query)).
						Str("group", "required by origin").
						Msg("Started")

					// Run the engine.
					errs, scraped := enginer.Search(query, opts, resChan)

					if len(errs) > 0 {
						searchOnce[engName].Errored()
						log.Error().
							Errs("errors", errs).
							Str("engine", engName.String()).
							Str("query", anonymize.String(query)).
							Str("group", "required by origin").
							Msg("Error searching")
					}

					if !scraped {
						log.Debug().
							Str("engine", engName.String()).
							Str("query", anonymize.String(query)).
							Str("group", "required by origin").
							Msg("Failed to scrape any results (probably timed out)")
					} else {
						searchOnce[engName].Scraped()
					}
				})

				// Indicate that the worker was successful.
				if searchOnce[engName].Success() {
					successOrigin.L.Lock()
					successOrigin.Signal()
					successOrigin.L.Unlock()
				}
			}()
		}
	}
}
