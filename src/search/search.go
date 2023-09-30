package search

import (
	"context"
	"net/url"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/sourcegraph/conc"
	"github.com/tminaorg/brzaguza/src/bucket"
	"github.com/tminaorg/brzaguza/src/bucket/result"
	"github.com/tminaorg/brzaguza/src/config"
	"github.com/tminaorg/brzaguza/src/engines"
	"github.com/tminaorg/brzaguza/src/rank"
)

func PerformSearch(query string, options engines.Options, config *config.Config) []result.Result {
	relay := bucket.Relay{
		ResultMap: make(map[string]*result.Result),
	}

	query = url.QueryEscape(query)

	resTiming := time.Now()
	log.Debug().Msg("Waiting for results from engines...")
	var worker conc.WaitGroup
	runEngines(config.Engines, query, &worker, &relay, options)
	worker.Wait()
	log.Debug().Msgf("Got results in %v", time.Since(resTiming).Milliseconds())

	rankTiming := time.Now()
	log.Debug().Msg("Ranking...")
	results := rank.Rank(relay.ResultMap)
	log.Debug().Msgf("Finished ranking in %v", time.Since(rankTiming).Milliseconds())

	log.Debug().Msg("Search Done!")

	return results
}

// engine_searcher, NewEngineStarter()  use this.
type EngineSearch func(context.Context, string, *bucket.Relay, engines.Options, config.Settings) error

func runEngines(engineMap map[string]config.Engine, query string, worker *conc.WaitGroup, relay *bucket.Relay, options engines.Options) {
	log.Info().Msgf("Enabled engines: %v", config.EnabledEngines)

	engineStarter := NewEngineStarter()
	for name, engine := range engineMap {
		engineName, nameErr := engines.NameString(name)
		if nameErr != nil {
			log.Panic().Err(nameErr).Msg("failed converting string to engine name")
			return
		}

		worker.Go(func() {
			if err := engineStarter[engineName](context.Background(), query, relay, options, engine.Settings); err != nil {
				log.Error().Err(err).Msgf("failed searching %v", engineName)
			}
		})
	}
}
