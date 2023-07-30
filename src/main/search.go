package main

import (
	"context"
	"sort"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/sourcegraph/conc"
	"github.com/tminaorg/brzaguza/src/engines/google"
	"github.com/tminaorg/brzaguza/src/structures"
)

var googleOptions structures.Options = structures.Options{
	UserAgent:  "",
	MaxPages:   3,
	ProxyAddr:  "",
	VisitPages: true,
}

func cleanQuery(query string) string {
	return strings.Replace(strings.Trim(query, " "), " ", "+", -1)
}

func performSearch(query string) []structures.Result {
	relay := structures.Relay{
		ResultChannel:     make(chan structures.Result),
		RankChannel:       make(chan structures.ResultRank),
		EngineDoneChannel: make(chan bool),
		ResultMap:         make(map[string]*structures.Result),
	}

	query = cleanQuery(query)

	const numberOfEngines int = 1
	var receivedEngines int = 0
	var worker conc.WaitGroup

	worker.Go(func() {
		err := google.Search(context.Background(), query, &relay, &googleOptions)
		if err != nil {
			log.Error().Err(err).Msg("Failed searching google.com")
		}
	})

	for receivedEngines < numberOfEngines {
		select {
		case result := <-relay.ResultChannel:
			mapRes, exists := relay.ResultMap[result.URL]
			if exists {
				if mapRes.Title == "" { // if rank was set first
					mapRes.Title = result.Title
				}
				mapRes.Description = result.Description // if rank was set first, or longer desc was found
			} else {
				relay.ResultMap[result.URL] = &result
			}
			log.Debug().Msgf("Got URL: %s", result.URL)
		case resRank := <-relay.RankChannel:
			mapRes, exists := relay.ResultMap[resRank.URL]
			if !exists { //if ResultRank came through channel before the Result
				relay.ResultMap[resRank.URL] = &structures.Result{
					Title:       "",
					Description: "",
					Rank:        resRank.Rank,
					URL:         resRank.URL,
				}
			} else {
				mapRes.Rank = resRank.Rank
			}

			log.Debug().Msgf("Updated rank to %d for %s: %s", resRank.Rank, mapRes.Title, resRank.URL)
		case <-relay.EngineDoneChannel:
			receivedEngines++
		}
	}

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
