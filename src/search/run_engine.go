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

func runSearcher(groupName string, onceWrap *onceWrapper, concMap *result.ResultConcMap, engName engines.Name, searcher scraper.Searcher, query string, opts options.Options) {
	// Run the engine only once.
	onceWrap.Do(func() {
		// Create a buffered channel for the results.
		resChan := make(chan result.ResultScraped, 100)

		// Start the receiver for the engine.
		var receiver sync.WaitGroup
		receiver.Add(1)
		go createReceiver(&receiver, resChan, concMap)

		log.Trace().
			Str("engine", engName.String()).
			Str("query", anonymize.String(query)).
			Str("group", groupName).
			Msg("Started")

		// Run the engine.
		errs, scraped := searcher.Search(query, opts, resChan)

		if len(errs) > 0 {
			onceWrap.Errored()
			log.Error().
				Errs("errors", errs).
				Str("engine", engName.String()).
				Str("query", anonymize.String(query)).
				Str("group", groupName).
				Msg("Error searching")
		}

		if !scraped {
			log.Debug().
				Str("engine", engName.String()).
				Str("query", anonymize.String(query)).
				Str("group", groupName).
				Msg("Failed to scrape any results (probably timed out)")
		} else {
			onceWrap.Scraped()
		}

		// Wait for the receiver to finish.
		receiver.Wait()
	})
}

func runSuggester(groupName string, onceWrap *onceWrapper, concMap *result.SuggestionConcMap, engName engines.Name, suggester scraper.Suggester, query string, locale options.Locale) {
	// Run the engine only once.
	onceWrap.Do(func() {
		// Create a buffered channel for the results.
		resChan := make(chan result.SuggestionScraped, 100)

		// Start the receiver for the engine.
		var receiver sync.WaitGroup
		receiver.Add(1)
		go createReceiver(&receiver, resChan, concMap)

		log.Trace().
			Str("engine", engName.String()).
			Str("query", anonymize.String(query)).
			Str("group", groupName).
			Msg("Started")

		// Run the engine.
		errs, scraped := suggester.Suggest(query, locale, resChan)

		if len(errs) > 0 {
			onceWrap.Errored()
			log.Error().
				Errs("errors", errs).
				Str("engine", engName.String()).
				Str("query", anonymize.String(query)).
				Str("group", groupName).
				Msg("Error searching")
		}

		if !scraped {
			log.Debug().
				Str("engine", engName.String()).
				Str("query", anonymize.String(query)).
				Str("group", groupName).
				Msg("Failed to scrape any results (probably timed out)")
		} else {
			onceWrap.Scraped()
		}

		// Wait for the receiver to finish.
		receiver.Wait()
	})
}
