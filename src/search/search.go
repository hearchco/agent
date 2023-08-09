package search

import (
	"context"
	"net/url"
	"sort"

	"github.com/rs/zerolog/log"
	"github.com/sourcegraph/conc"
	"github.com/tminaorg/brzaguza/src/config"
	"github.com/tminaorg/brzaguza/src/engines/bing"
	"github.com/tminaorg/brzaguza/src/engines/brave"
	"github.com/tminaorg/brzaguza/src/engines/duckduckgo"
	"github.com/tminaorg/brzaguza/src/engines/etools"
	"github.com/tminaorg/brzaguza/src/engines/google"
	"github.com/tminaorg/brzaguza/src/engines/mojeek"
	"github.com/tminaorg/brzaguza/src/engines/qwant"
	"github.com/tminaorg/brzaguza/src/engines/startpage"
	"github.com/tminaorg/brzaguza/src/engines/swisscows"
	"github.com/tminaorg/brzaguza/src/structures"
)

func PerformSearch(query string, maxPages int, visitPages bool, config *config.Config) []structures.Result {
	relay := structures.Relay{
		ResultMap:         make(map[string]*structures.Result),
		EngineDoneChannel: make(chan bool),
	}

	options := structures.SEOptions{
		MaxPages:   maxPages,
		VisitPages: visitPages,
	}

	query = url.QueryEscape(query)

	var worker conc.WaitGroup
	runEngines(config.Engines, query, &worker, &relay, &options)
	worker.Wait()

	var results []structures.Result = make([]structures.Result, 0, len(relay.ResultMap))
	for _, res := range relay.ResultMap {
		results = append(results, *res)
	}

	sort.Sort(structures.ByRank(results))

	log.Debug().Msg("All processing done, waiting for closing of goroutines.")
	worker.Wait()

	log.Debug().Msg("Done! Received All Engines!")

	return results
}

func runEngines(engines []config.Engine, query string, worker *conc.WaitGroup, relay *structures.Relay, options *structures.SEOptions) {
	log.Info().Msgf("Enabled engines: %v", config.EnabledEngines)

	for _, engine := range engines {
		switch engine.Name {
		case structures.Google:
			worker.Go(func() {
				err := google.Search(context.Background(), query, relay, options, &engine.Settings)
				if err != nil {
					log.Error().Err(err).Msgf("Failed searching %v", google.SEDomain)
				}
			})
		case structures.DuckDuckGo:
			worker.Go(func() {
				err := duckduckgo.Search(context.Background(), query, relay, options, &engine.Settings)
				if err != nil {
					log.Error().Err(err).Msgf("Failed searching %v", duckduckgo.SEDomain)
				}
			})
		case structures.Mojeek:
			worker.Go(func() {
				err := mojeek.Search(context.Background(), query, relay, options, &engine.Settings)
				if err != nil {
					log.Error().Err(err).Msgf("Failed searching %v", mojeek.SEDomain)
				}
			})
		case structures.Qwant:
			worker.Go(func() {
				err := qwant.Search(context.Background(), query, relay, options, &engine.Settings)
				if err != nil {
					log.Error().Err(err).Msgf("Failed searching %v", qwant.SEDomain)
				}
			})
		case structures.Etools:
			worker.Go(func() {
				err := etools.Search(context.Background(), query, relay, options, &engine.Settings)
				if err != nil {
					log.Error().Err(err).Msgf("Failed searching %v", etools.SEDomain)
				}
			})
		case structures.Swisscows:
			worker.Go(func() {
				err := swisscows.Search(context.Background(), query, relay, options, &engine.Settings)
				if err != nil {
					log.Error().Err(err).Msgf("Failed searching %v", swisscows.SEDomain)
				}
			})
		case structures.Brave:
			worker.Go(func() {
				err := brave.Search(context.Background(), query, relay, options, &engine.Settings)
				if err != nil {
					log.Error().Err(err).Msgf("Failed searching %v", brave.SEDomain)
				}
			})
		case structures.Bing:
			worker.Go(func() {
				err := bing.Search(context.Background(), query, relay, options, &engine.Settings)
				if err != nil {
					log.Error().Err(err).Msgf("Failed searching %v", bing.SEDomain)
				}
			})
		case structures.Startpage:
			worker.Go(func() {
				err := startpage.Search(context.Background(), query, relay, options, &engine.Settings)
				if err != nil {
					log.Error().Err(err).Msgf("Failed searching %v", startpage.SEDomain)
				}
			})
		}
	}
}
