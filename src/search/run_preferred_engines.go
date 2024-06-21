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

func runPreferredEngines(searchers []scraper.Searcher, wgPreferredEngines *sync.WaitGroup, query string, opts options.Options, preferredEngines []engines.Name, engChan chan chan result.ResultScraped, searchOnce map[engines.Name]*onceWrapper) {
	wgPreferredEngines.Add(len(preferredEngines))
	for _, engName := range preferredEngines {
		searcher := searchers[engName]
		resChan := make(chan result.ResultScraped, 100)
		engChan <- resChan
		go func() {
			defer wgPreferredEngines.Done()
			searchOnce[engName].Do(func() {
				log.Trace().
					Str("engine", engName.String()).
					Str("query", anonymize.String(query)).
					Str("group", "preferred").
					Msg("Started")

				// Run the engine.
				errs, scraped := searcher.Search(query, opts, resChan)

				if len(errs) > 0 {
					searchOnce[engName].Errored()
					log.Error().
						Errs("errors", errs).
						Str("engine", engName.String()).
						Str("query", anonymize.String(query)).
						Str("group", "preferred").
						Msg("Error searching")
				}

				if !scraped {
					log.Debug().
						Str("engine", engName.String()).
						Str("query", anonymize.String(query)).
						Str("group", "preferred").
						Msg("Failed to scrape any results (probably timed out)")
				} else {
					searchOnce[engName].Scraped()
				}
			})
		}()
	}
}
