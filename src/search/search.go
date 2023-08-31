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
	"github.com/tminaorg/brzaguza/src/engines/bing"
	"github.com/tminaorg/brzaguza/src/engines/brave"
	"github.com/tminaorg/brzaguza/src/engines/duckduckgo"
	"github.com/tminaorg/brzaguza/src/engines/etools"
	"github.com/tminaorg/brzaguza/src/engines/google"
	"github.com/tminaorg/brzaguza/src/engines/mojeek"
	"github.com/tminaorg/brzaguza/src/engines/presearch"
	"github.com/tminaorg/brzaguza/src/engines/qwant"
	"github.com/tminaorg/brzaguza/src/engines/startpage"
	"github.com/tminaorg/brzaguza/src/engines/swisscows"
	"github.com/tminaorg/brzaguza/src/engines/yep"
	"github.com/tminaorg/brzaguza/src/rank"
)

func PerformSearch(query string, maxPages int, visitPages bool, config *config.Config) []result.Result {
	relay := bucket.Relay{
		ResultMap:         make(map[string]*result.Result),
		EngineDoneChannel: make(chan bool),
	}

	options := engines.Options{
		MaxPages:   maxPages,
		VisitPages: visitPages,
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

func runEngines(engineMap map[engines.Name]config.Engine, query string, worker *conc.WaitGroup, relay *bucket.Relay, options engines.Options) {
	log.Info().Msgf("Enabled engines: %v", config.EnabledEngines)

	for name, engine := range engineMap {
		switch name {
		case engines.Google:
			worker.Go(func() {
				err := google.Search(context.Background(), query, relay, options, engine.Settings)
				if err != nil {
					log.Error().Err(err).Msgf("Failed searching %v", google.Info.Domain)
				}
			})
		case engines.DuckDuckGo:
			worker.Go(func() {
				err := duckduckgo.Search(context.Background(), query, relay, options, engine.Settings)
				if err != nil {
					log.Error().Err(err).Msgf("Failed searching %v", duckduckgo.Info.Domain)
				}
			})
		case engines.Mojeek:
			worker.Go(func() {
				err := mojeek.Search(context.Background(), query, relay, options, engine.Settings)
				if err != nil {
					log.Error().Err(err).Msgf("Failed searching %v", mojeek.Info.Domain)
				}
			})
		case engines.Qwant:
			worker.Go(func() {
				err := qwant.Search(context.Background(), query, relay, options, engine.Settings)
				if err != nil {
					log.Error().Err(err).Msgf("Failed searching %v", qwant.Info.Domain)
				}
			})
		case engines.Etools:
			worker.Go(func() {
				err := etools.Search(context.Background(), query, relay, options, engine.Settings)
				if err != nil {
					log.Error().Err(err).Msgf("Failed searching %v", etools.Info.Domain)
				}
			})
		case engines.Swisscows:
			worker.Go(func() {
				err := swisscows.Search(context.Background(), query, relay, options, engine.Settings)
				if err != nil {
					log.Error().Err(err).Msgf("Failed searching %v", swisscows.Info.Domain)
				}
			})
		case engines.Brave:
			worker.Go(func() {
				err := brave.Search(context.Background(), query, relay, options, engine.Settings)
				if err != nil {
					log.Error().Err(err).Msgf("Failed searching %v", brave.Info.Domain)
				}
			})
		case engines.Bing:
			worker.Go(func() {
				err := bing.Search(context.Background(), query, relay, options, engine.Settings)
				if err != nil {
					log.Error().Err(err).Msgf("Failed searching %v", bing.Info.Domain)
				}
			})
		case engines.Startpage:
			worker.Go(func() {
				err := startpage.Search(context.Background(), query, relay, options, engine.Settings)
				if err != nil {
					log.Error().Err(err).Msgf("Failed searching %v", startpage.Info.Domain)
				}
			})
		case engines.Yep:
			worker.Go(func() {
				err := yep.Search(context.Background(), query, relay, options, engine.Settings)
				if err != nil {
					log.Error().Err(err).Msgf("Failed searching %v", yep.Info.Domain)
				}
			})
		case engines.Presearch:
			worker.Go(func() {
				err := presearch.Search(context.Background(), query, relay, options, engine.Settings)
				if err != nil {
					log.Error().Err(err).Msgf("Failed searching %v", presearch.Info.Domain)
				}
			})
		}
	}
}
