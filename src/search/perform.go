package search

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/hearchco/hearchco/src/anonymize"
	"github.com/hearchco/hearchco/src/config"
	"github.com/hearchco/hearchco/src/search/bucket"
	"github.com/hearchco/hearchco/src/search/engines"
	"github.com/hearchco/hearchco/src/search/rank"
	"github.com/hearchco/hearchco/src/search/result"
	"github.com/rs/zerolog/log"
)

// engine_searcher -> NewEngineStarter() uses this
type EngineSearch func(context.Context, string, *bucket.Relay, engines.Options, config.Settings, config.Timings) []error

func PerformSearch(query string, options engines.Options, conf config.Config) []result.Result {
	searchTimer := time.Now()

	query, timings, enginesToRun := procBang(query, &options, conf)

	query = url.QueryEscape(query)
	log.Debug().
		Str("queryAnon", anonymize.String(query)).
		Str("queryHash", anonymize.HashToSHA256B64(query)).
		Msg("Searching")

	resTimer := time.Now()
	log.Debug().Msg("Waiting for results from engines...")
	relay := runEngines(enginesToRun, query, options, conf.Settings, timings)
	log.Debug().
		Int64("ms", time.Since(resTimer).Milliseconds()).
		Msg("Got results")

	rankTimer := time.Now()
	log.Debug().Msg("Ranking...")
	results := rank.Rank(relay.ResultMap, conf.Categories[options.Category].Ranking) // have to make copy, since its a map value
	rankTimeSince := time.Since(rankTimer)
	log.Debug().
		Int64("ms", rankTimeSince.Milliseconds()).
		Int64("ns", rankTimeSince.Nanoseconds()).
		Msg("Finished ranking")

	log.Debug().
		Int64("ms", time.Since(searchTimer).Milliseconds()).
		Msg("Found results")

	return results
}

func runEngines(engs []engines.Name, query string, options engines.Options, settings map[engines.Name]config.Settings, timings config.Timings) *bucket.Relay {
	config.EnabledEngines = engs
	log.Info().
		Int("number", len(config.EnabledEngines)).
		Str("engines", fmt.Sprintf("%v", config.EnabledEngines)).
		Msg("Enabled engines")

	relay := &bucket.Relay{
		ResultMap: make(map[string]*result.Result),
	}

	var wg sync.WaitGroup
	engineStarter := NewEngineStarter()

	for i := range engs {
		wg.Add(1)
		eng := engs[i] // dont change for to `for _, eng := range engs {`, eng retains the same address throughout the whole loop
		go func() {
			defer wg.Done()
			// if an error can be handled inside, it wont be returned
			// runs the Search function in the engine package
			errs := engineStarter[eng](context.Background(), query, relay, options, settings[eng], timings)
			if len(errs) > 0 {
				log.Error().
					Errs("errors", errs).
					Str("engine", eng.String()).
					Msg("search.runEngines(): error while searching")
			}
		}()
	}

	wg.Wait()
	return relay
}
