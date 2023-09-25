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
		ResultMap:         make(map[string]*result.Result),
		EngineDoneChannel: make(chan bool),
	}

	query = url.QueryEscape(query)

	var worker conc.WaitGroup
	runEngines(config.Engines, query, &worker, &relay, options)
	worker.Wait()

	results := make([]result.Result, 0, len(relay.ResultMap))
	for _, res := range relay.ResultMap {
		results = append(results, *res)
	}

	sort.Sort(rank.ByRank(results))

	log.Debug().Msg("All processing done, waiting for closing of goroutines.")
	worker.Wait()

	log.Debug().Msg("Done! Received All Engines!")

	return results
}

func runEngines(engineMap map[string]config.Engine, query string, worker *conc.WaitGroup, relay *bucket.Relay, options engines.Options) {
	log.Info().Msgf("Enabled engines: %v", config.EnabledEngines)

	engineStarter := NewEngineStarter()
	for name, engine := range engineMap {
		engineName, nameErr := engines.NameString(name)
		if nameErr != nil {
			log.Panic().Err(nameErr).Msg("failed converting string to engine name")
			return
		}

		err := engineStarter[engineName](context.Background(), query, relay, options, engine.Settings)
		if err != nil {
			log.Error().Err(err).Msgf("failed searching %v", engineName)
		}
	}
}
