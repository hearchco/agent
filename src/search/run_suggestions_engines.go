package search

import (
	"context"
	"sync"

	"github.com/rs/zerolog/log"

	"github.com/hearchco/agent/src/search/engines"
	"github.com/hearchco/agent/src/search/engines/options"
	"github.com/hearchco/agent/src/search/result"
	"github.com/hearchco/agent/src/search/scraper"
	"github.com/hearchco/agent/src/utils/anonymize"
)

func runSuggestionsEngines(suggesters []scraper.Suggester, cancelCtx context.CancelFunc, query string, locale options.Locale, enabledEngines []engines.Name, engChan chan chan []result.SuggestionScraped) {
	// Wait for any engine to successfully finish.
	c := sync.NewCond(&sync.Mutex{})
	go func() {
		c.L.Lock()
		c.Wait()
		c.L.Unlock()
		cancelCtx()
	}()

	// Wait for all engines to finish (successful or not).
	var wg sync.WaitGroup
	wg.Add(len(suggesters))
	go func() {
		wg.Wait()
		cancelCtx()
	}()

	// Run each engine.
	for _, engName := range enabledEngines {
		suggester := suggesters[engName]
		sugChan := make(chan []result.SuggestionScraped)
		engChan <- sugChan
		go func() {
			defer wg.Done()
			err, found := suggester.Suggest(query, locale, sugChan)
			if err != nil {
				log.Error().
					Err(err).
					Str("engine", suggester.GetName().String()).
					Str("query", anonymize.String(query)).
					Msg("Suggest failed")
			} else if !found {
				log.Error().
					Str("engine", suggester.GetName().String()).
					Str("query", anonymize.String(query)).
					Msg("No suggestions found")
			} else {
				c.L.Lock()
				c.Signal()
				c.L.Unlock()
			}
		}()
	}
}
