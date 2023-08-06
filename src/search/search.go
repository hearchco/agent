package search

import (
	"context"
	"net/url"
	"sort"

	"github.com/rs/zerolog/log"
	"github.com/sourcegraph/conc"
	"github.com/tminaorg/brzaguza/src/engines/duckduckgo"
	"github.com/tminaorg/brzaguza/src/engines/etools"
	"github.com/tminaorg/brzaguza/src/engines/google"
	"github.com/tminaorg/brzaguza/src/engines/mojeek"
	"github.com/tminaorg/brzaguza/src/engines/qwant"
	"github.com/tminaorg/brzaguza/src/engines/startpage"
	"github.com/tminaorg/brzaguza/src/structures"
)

type Engine int

const (
	Google Engine = iota
	Mojeek
	DuckDuckGo
	Qwant
	Etools
	Startpage
)

func PerformSearch(query string, maxPages int, visitPages bool) []structures.Result {
	relay := structures.Relay{
		ResultMap:         make(map[string]*structures.Result),
		EngineDoneChannel: make(chan bool),
	}

	options := structures.Options{
		MaxPages:   maxPages,
		VisitPages: visitPages,
	}

	query = url.QueryEscape(query)

	var worker conc.WaitGroup
	var toSearch []Engine = []Engine{Startpage}
	runEngines(toSearch, query, &worker, &relay, &options)
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

func runEngines(toSearch []Engine, query string, worker *conc.WaitGroup, relay *structures.Relay, options *structures.Options) {
	for _, eng := range toSearch {
		switch Engine(eng) {
		case Google:
			worker.Go(func() {
				err := google.Search(context.Background(), query, relay, options)
				if err != nil {
					log.Error().Err(err).Msgf("Failed searching %v", google.SEDomain)
				}
			})
		case DuckDuckGo:
			worker.Go(func() {
				err := duckduckgo.Search(context.Background(), query, relay, options)
				if err != nil {
					log.Error().Err(err).Msgf("Failed searching %v", duckduckgo.SEDomain)
				}
			})
		case Mojeek:
			worker.Go(func() {
				err := mojeek.Search(context.Background(), query, relay, options)
				if err != nil {
					log.Error().Err(err).Msgf("Failed searching %v", mojeek.SEDomain)
				}
			})
		case Qwant:
			worker.Go(func() {
				err := qwant.Search(context.Background(), query, relay, options)
				if err != nil {
					log.Error().Err(err).Msgf("Failed searching %v", qwant.SEDomain)
				}
			})
		case Etools:
			worker.Go(func() {
				err := etools.Search(context.Background(), query, relay, options)
				if err != nil {
					log.Error().Err(err).Msgf("Failed searching %v", etools.SEDomain)
				}
			})
		case Startpage:
			worker.Go(func() {
				err := startpage.Search(context.Background(), query, relay, options)
				if err != nil {
					log.Error().Err(err).Msgf("Failed searching %v", startpage.SEDomain)
				}
			})
		}
	}
}
