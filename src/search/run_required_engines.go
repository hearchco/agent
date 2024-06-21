package search

import (
	"sync"

	"github.com/rs/zerolog/log"

	"github.com/hearchco/agent/src/search/engines"
	"github.com/hearchco/agent/src/search/engines/options"
	"github.com/hearchco/agent/src/search/result"
	"github.com/hearchco/agent/src/search/scraper"
	"github.com/hearchco/agent/src/utils/anonymize"
)

func runRequiredEngines(searchers []scraper.Searcher, wgRequiredEngines *sync.WaitGroup, query string, opts options.Options, requiredEngines []engines.Name, engChan chan chan result.ResultScraped, searchOnce map[engines.Name]*onceWrapper) {
	wgRequiredEngines.Add(len(requiredEngines))
	for _, engName := range requiredEngines {
		searcher := searchers[engName]
		resChan := make(chan result.ResultScraped, 100)
		engChan <- resChan
		go func() {
			defer wgRequiredEngines.Done()
			searchOnce[engName].Do(func() {
				log.Trace().
					Str("engine", engName.String()).
					Str("query", anonymize.String(query)).
					Str("group", "required").
					Msg("Started")

				// Run the engine.
				errs, scraped := searcher.Search(query, opts, resChan)

				if len(errs) > 0 {
					searchOnce[engName].Errored()
					log.Error().
						Errs("errors", errs).
						Str("engine", engName.String()).
						Str("query", anonymize.String(query)).
						Str("group", "required").
						Msg("Error searching")
				}

				if !scraped {
					log.Debug().
						Str("engine", engName.String()).
						Str("query", anonymize.String(query)).
						Str("group", "required").
						Msg("Failed to scrape any results (probably timed out)")
				} else {
					searchOnce[engName].Scraped()
				}
			})
		}()
	}
}
