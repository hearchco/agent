package search

import (
	"context"
	"sort"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/sourcegraph/conc"
	"github.com/tminaorg/brzaguza/src/engines/duckduckgo"
	"github.com/tminaorg/brzaguza/src/engines/google"
	"github.com/tminaorg/brzaguza/src/engines/mojeek"
	"github.com/tminaorg/brzaguza/src/structures"
)

func cleanQuery(query string) string {
	return strings.Replace(strings.Trim(query, " "), " ", "+", -1)
}

func PerformSearch(query string, maxPages int, visitPages bool) []structures.Result {
	relay := structures.Relay{
		ResultMap:         make(map[string]*structures.Result),
		EngineDoneChannel: make(chan bool),
	}

	options := structures.Options{
		MaxPages:   maxPages,
		VisitPages: visitPages,
	}

	query = cleanQuery(query)

	var worker conc.WaitGroup

	worker.Go(func() {
		err := google.Search(context.Background(), query, &relay, &options)
		if err != nil {
			log.Error().Err(err).Msg("Failed searching google.com")
		}
	})

	worker.Go(func() {
		err := duckduckgo.Search(context.Background(), query, &relay, &options)
		if err != nil {
			log.Error().Err(err).Msg("Failed searching lite.duckduckgo.com")
		}
	})

	worker.Go(func() {
		err := mojeek.Search(context.Background(), query, &relay, &options)
		if err != nil {
			log.Error().Err(err).Msg("Failed searching mojeek.com")
		}
	})

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
