package search

import (
	"context"
	"net/url"
	"strings"
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
	runEngines(config.Engines, query, &worker, &relay, options)
	worker.Wait()
	log.Debug().Msgf("Got results in %vms", time.Since(resTimer).Milliseconds())

	rankTimer := time.Now()
	log.Debug().Msg("Ranking...")
	results := rank.Rank(relay.ResultMap, &(config.Ranking))
	log.Debug().Msgf("Finished ranking in %vns", time.Since(rankTimer).Nanoseconds())

	log.Debug().Msgf("Found results in %vms", time.Since(searchTimer).Milliseconds())

	return results
}

// engine_searcher, NewEngineStarter()  use this.
type EngineSearch func(context.Context, string, *bucket.Relay, engines.Options, config.Settings) error

func runEngines(engineMap map[string]config.Engine, query string, worker *conc.WaitGroup, relay *bucket.Relay, options engines.Options) {
	log.Info().Msgf("Enabled engines (%v): %v", len(config.EnabledEngines), config.EnabledEngines)

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

func setCategory(query string, options *engines.Options) {
	category := extractCategory(query)
	if category != "" {
		options.Category = category
	}
}

func extractCategory(query string) string {
	valid := []string{"info", "science", "news", "blog", "surf", "newnews", "wiki", "sci", "nnews"}
	// info[/wiki], science[/sci], newnews[/nnews]

	if query[0] != '!' {
		return ""
	}
	cat := strings.SplitN(query, " ", 2)[0][1:]

	ok := false
	for i := range valid {
		if valid[i] == cat {
			ok = true
			break
		}
	}

	if ok {
		return cat
	}
	return ""
}
