package search

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/hearchco/hearchco/src/anonymize"
	"github.com/hearchco/hearchco/src/config"
	"github.com/hearchco/hearchco/src/search/bucket"
	"github.com/hearchco/hearchco/src/search/engines"
	"github.com/hearchco/hearchco/src/search/rank"
	"github.com/hearchco/hearchco/src/search/result"
	"github.com/rs/zerolog/log"
	"github.com/sourcegraph/conc"
)

// engine_searcher -> NewEngineStarter() uses this
type EngineSearch func(context.Context, string, *bucket.Relay, engines.Options, config.Settings, config.Timings) error

func PerformSearch(query string, options engines.Options, conf *config.Config) []result.Result {
	searchTimer := time.Now()

	relay := bucket.Relay{
		ResultMap: make(map[string]*result.Result),
	}

	query, timings, toRun := procBang(query, &options, conf)

	query = url.QueryEscape(query)
	log.Debug().
		Str("queryAnon", anonymize.String(query)).
		Str("queryHash", anonymize.HashToSHA256B64(query)).
		Msg("Searching")

	resTimer := time.Now()
	log.Debug().Msg("Waiting for results from engines...")
	var worker conc.WaitGroup
	runEngines(toRun, timings, conf.Settings, query, &worker, &relay, options)
	worker.Wait()
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

func runEngines(engs []engines.Name, timings config.Timings, settings map[engines.Name]config.Settings, query string, worker *conc.WaitGroup, relay *bucket.Relay, options engines.Options) {
	config.EnabledEngines = engs
	log.Info().
		Int("number", len(config.EnabledEngines)).
		Str("engines", fmt.Sprintf("%v", config.EnabledEngines)).
		Msg("Enabled engines")

	engineStarter := NewEngineStarter()
	for i := range engs {
		eng := engs[i] // dont change for to `for _, eng := range engs {`, eng retains the same address throughout the whole loop
		worker.Go(func() {
			// if an error can be handled inside, it wont be returned
			// runs the Search function in the engine package
			err := engineStarter[eng](context.Background(), query, relay, options, settings[eng], timings)
			if err != nil {
				log.Error().
					Err(err).
					Str("engine", eng.String()).
					Msg("search.runEngines(): error while searching")
			}
		})
	}
}
