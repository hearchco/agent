package search

import (
	"context"
	"sort"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/sourcegraph/conc"
	"github.com/tminaorg/brzaguza/src/engines/google"
	"github.com/tminaorg/brzaguza/src/rank"
	"github.com/tminaorg/brzaguza/src/structures"
)

func cleanQuery(query string) string {
	return strings.Replace(strings.Trim(query, " "), " ", "+", -1)
}

func PerformSearch(query string, maxPages int, visitPages bool) []structures.Result {
	relay := structures.Relay{
		ResultChannel:     make(chan structures.Result),
		ResponseChannel:   make(chan structures.ResultResponse),
		EngineDoneChannel: make(chan bool),
		ResultMap:         make(map[string]*structures.Result),
	}

	options := structures.Options{
		MaxPages:   maxPages,
		VisitPages: visitPages,
	}

	query = cleanQuery(query)

	const numberOfEngines int = 1
	var receivedEngines int = 0
	var worker conc.WaitGroup

	worker.Go(func() {
		err := google.Search(context.Background(), query, &relay, &options)
		if err != nil {
			log.Error().Err(err).Msg("Failed searching google.com")
		}
	})

	for receivedEngines < numberOfEngines {
		select {
		case result := <-relay.ResultChannel:
			log.Trace().Msgf("Got URL: %v", result.URL)

			mapRes, exists := relay.ResultMap[result.URL]
			if exists {
				if mapRes.Title == "" { // if response was set first
					mapRes.Title = result.Title
					mapRes.SEPage = result.SEPage
					mapRes.SEPageRank = result.SEPageRank
					rank.SetRank(mapRes)
				}
				mapRes.Description = result.Description // if response was set first, or longer desc was found
			} else {
				relay.ResultMap[result.URL] = &result
			}
		case resRes := <-relay.ResponseChannel:
			log.Trace().Msgf("Got response for %v", resRes.URL)

			mapRes, exists := relay.ResultMap[resRes.URL]
			if exists {
				mapRes.Response = resRes.Response
				rank.SetRank(mapRes)
			} else {
				//if ResultRank came through channel before the Result
				relay.ResultMap[resRes.URL] = &structures.Result{
					Response: resRes.Response,
					Rank:     -1,
				}
			}
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
