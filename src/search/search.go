package search

import (
	"context"
	"net/url"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/sourcegraph/conc"
	"github.com/tminaorg/brzaguza/src/bucket"
	"github.com/tminaorg/brzaguza/src/bucket/result"
	"github.com/tminaorg/brzaguza/src/category"
	"github.com/tminaorg/brzaguza/src/config"
	"github.com/tminaorg/brzaguza/src/engines"
	"github.com/tminaorg/brzaguza/src/rank"
)

func PerformSearch(query string, options engines.Options, conf *config.Config) []result.Result {
	searchTimer := time.Now()

	relay := bucket.Relay{
		ResultMap: make(map[string]*result.Result),
	}

	setCategory(query, &options)

	query = url.QueryEscape(query)
	log.Debug().Msg(query)

	resTimer := time.Now()
	log.Debug().Msg("Waiting for results from engines...")
	var worker conc.WaitGroup
	runEngines(conf.Categories[options.Category].Engines, conf.Settings, query, &worker, &relay, options)
	worker.Wait()
	log.Debug().Msgf("Got results in %vms", time.Since(resTimer).Milliseconds())

	rankTimer := time.Now()
	log.Debug().Msg("Ranking...")
	results := rank.Rank(relay.ResultMap, conf.Categories[options.Category].Ranking) // have to make copy, since its a map value
	log.Debug().Msgf("Finished ranking in %vns", time.Since(rankTimer).Nanoseconds())

	log.Debug().Msgf("Found results in %vms", time.Since(searchTimer).Milliseconds())

	return results
}

// engine_searcher, NewEngineStarter()  use this.
type EngineSearch func(context.Context, string, *bucket.Relay, engines.Options, config.Settings) error

func runEngines(engineMap map[string]config.Engine, settings map[string]config.Settings, query string, worker *conc.WaitGroup, relay *bucket.Relay, options engines.Options) {
	engineStarter := NewEngineStarter()
	for name, engine := range engineMap {
		if !engine.Enabled {
			continue
		}

		engineName, nameErr := engines.NameString(name)
		if nameErr != nil {
			log.Panic().Err(nameErr).Msg("failed converting string to engine name")
			return
		}

		config.EnabledEngines = append(config.EnabledEngines, engineName)

		worker.Go(func() {
			if err := engineStarter[engineName](context.Background(), query, relay, options, settings[name]); err != nil {
				log.Error().Err(err).Msgf("failed searching %v", engineName)
			}
		})
	}

	log.Info().Msgf("Enabled engines (%v): %v", len(config.EnabledEngines), config.EnabledEngines)
}

func setCategory(query string, options *engines.Options) {
	cat := category.FromQuery(query)
	if cat != "" {
		options.Category = cat
	}
	if options.Category == "" {
		options.Category = category.GENERAL
	}
}
