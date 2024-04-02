package search

import (
	"context"
	"net/url"
	"sync"
	"time"

	"github.com/hearchco/hearchco/src/anonymize"
	"github.com/hearchco/hearchco/src/config"
	"github.com/hearchco/hearchco/src/search/bucket"
	"github.com/hearchco/hearchco/src/search/category"
	"github.com/hearchco/hearchco/src/search/engines"
	"github.com/hearchco/hearchco/src/search/rank"
	"github.com/hearchco/hearchco/src/search/result"
	"github.com/rs/zerolog/log"
)

// engine_searcher -> NewEngineStarter() uses this
type EngineSearch func(context.Context, string, *bucket.Relay, engines.Options, config.Settings, config.Timings, string, int) []error

func PerformSearch(query string, options engines.Options, settings map[engines.Name]config.Settings, categories map[category.Name]config.Category, salt string) []result.Result {
	if query == "" {
		log.Trace().Msg("Empty search query.")
		return []result.Result{}
	}

	searchTimer := time.Now()

	query, cat, timings, enginesToRun := procBang(query, options.Category, settings, categories)
	// set the new category only within the scope of this function
	options.Category = cat
	query = url.QueryEscape(query)

	// check again after the bang is taken out
	if query == "" {
		log.Trace().Msg("Empty search query (with bang present).")
		return []result.Result{}
	}

	log.Debug().
		Str("queryAnon", anonymize.String(query)).
		Str("queryHash", anonymize.HashToSHA256B64(query)).
		Msg("Searching")

	resTimer := time.Now()
	log.Debug().Msg("Waiting for results from engines...")

	resultMap := runEngines(enginesToRun, query, options, settings, timings, salt)

	log.Debug().
		Dur("duration", time.Since(resTimer)).
		Msg("Got results")

	rankTimer := time.Now()
	log.Debug().Msg("Ranking...")

	results := rank.Rank(resultMap, categories[options.Category].Ranking)

	log.Debug().
		Dur("duration", time.Since(rankTimer)).
		Msg("Finished ranking")

	log.Debug().
		Dur("duration", time.Since(searchTimer)).
		Msg("Found results")

	return results
}

func runEngines(engs []engines.Name, query string, options engines.Options, settings map[engines.Name]config.Settings, timings config.Timings, salt string) map[string]*result.Result {
	engsStrs := make([]string, 0, len(engs))
	for _, eng := range engs {
		engsStrs = append(engsStrs, eng.String())
	}

	log.Info().
		Int("number", len(engs)).
		Strs("engines", engsStrs).
		Msg("Enabled engines")

	relay := bucket.Relay{
		ResultMap: make(map[string]*result.Result),
	}

	var wg sync.WaitGroup
	engineStarter := NewEngineStarter()

	for _, eng := range engs {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// if an error can be handled inside, it won't be returned
			// runs the Search function in the engine package
			errs := engineStarter[eng](context.Background(), query, &relay, options, settings[eng], timings, salt, len(engs))
			if len(errs) > 0 {
				log.Error().
					Errs("errors", errs).
					Str("engine", eng.String()).
					Msg("search.runEngines(): error while searching")
			}
		}()
	}

	wg.Wait()
	return relay.ResultMap
}
