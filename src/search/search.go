package search

import (
	"context"
	"net/url"
	"sort"

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

	var worker conc.WaitGroup
	runEngines(config.Engines, query, &worker, &relay, options)
	log.Debug().Msg("Waiting for results from engines...")
	worker.Wait()

	results := make([]result.Result, 0, len(relay.ResultMap))
	for _, res := range relay.ResultMap {
		results = append(results, *res)
	}

	sort.Sort(rank.ByRank(results))

	log.Debug().Msg("All processing done!")

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
