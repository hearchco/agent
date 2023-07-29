package main

import (
	"context"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/sourcegraph/conc"
	google "github.com/tminaorg/brzaguza/src/engines"
	"github.com/tminaorg/brzaguza/src/relay"
	"github.com/tminaorg/brzaguza/src/structures"
)

var googleOptions structures.Options = structures.Options{
	UserAgent:     "",
	Limit:         50,
	ProxyAddr:     "",
	JustFirstPage: false,
}

func performSearch(query string) []structures.Result {
	relay.ResultMap = make(map[string]*structures.Result) // probably faster than clearing current map

	query = strings.Trim(query, " ")
	query = strings.Replace(query, " ", "+", -1)

	const numberOfEngines int = 1
	var receivedEngines int = 0
	var worker conc.WaitGroup

	worker.Go(func() {
		err := google.Search(context.Background(), query, &googleOptions, &worker)
		if err != nil {
			log.Error().Err(err).Msg("Failed searching google.com")
		}
	})

	for receivedEngines < numberOfEngines {
		select {
		case result := <-relay.ResultChannel:
			relay.ResultMap[result.URL] = &result
			log.Debug().Msgf("Got URL: %v\n", result.URL)
		case resRank := <-relay.RankChannel:
			relay.ResultMap[resRank.URL].Rank = resRank.Rank
			log.Debug().Msgf("Updated rank to %v for %v\n", resRank.Rank, resRank.URL)
		case <-relay.EngineDoneChannel:
			receivedEngines++
		}
	}

	var results []structures.Result = make([]structures.Result, 0, len(relay.ResultMap))
	for _, res := range relay.ResultMap {
		results = append(results, *res)
	}

	log.Debug().Msg("All processing done, waiting for closing of goroutines.")
	worker.Wait()

	log.Debug().Msg("Done! Received All Engines!")

	return results
}
