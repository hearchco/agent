package search

import (
	"github.com/rs/zerolog/log"

	"github.com/hearchco/agent/src/search/engines"
	"github.com/hearchco/agent/src/search/result"
)

func createReceiver(engChan chan chan result.ResultScraped, results *result.ConcurrentMap, enabledEnginesLen int) {
	for resChan := range engChan {
		go func() {
			for recRes := range resChan {
				if recRes.Rank().SearchEngine().String() == "" || recRes.Rank().SearchEngine() == engines.UNDEFINED {
					log.Panic().
						Str("engine", recRes.Rank().SearchEngine().String()).
						Msg("Received a result with an undefined search engine")
					// ^PANIC - Assert because it should never happen.
				}

				// Lock the results map due to modifications.
				results.Mutex.Lock()

				mapRes, exists := results.Map[recRes.URL()]
				if !exists {
					// Add the result to the results map.
					results.Map[recRes.URL()] = recRes.Convert(enabledEnginesLen)
				} else {
					var alreadyIn *result.Rank

					// Check if the engine rank is already in the result.
					for i, er := range mapRes.EngineRanks() {
						if recRes.Rank().SearchEngine() == er.SearchEngine() {
							alreadyIn = &mapRes.EngineRanks()[i]
							break
						}
					}

					// Update the result if the new rank is better.
					if alreadyIn == nil {
						mapRes.AppendEngineRanks(recRes.Rank().Convert())
					} else if alreadyIn.Page() > recRes.Rank().Page() {
						alreadyIn.SetPage(recRes.Rank().Page(), recRes.Rank().OnPageRank())
					} else if alreadyIn.Page() == recRes.Rank().Page() && alreadyIn.OnPageRank() > recRes.Rank().OnPageRank() {
						alreadyIn.SetOnPageRank(recRes.Rank().OnPageRank())
					}

					// Update the description if the new description is longer.
					if len(mapRes.Description()) < len(recRes.Description()) {
						mapRes.SetDescription(recRes.Description())
					}
				}

				// Unlock the results map.
				results.Mutex.Unlock()
			}
		}()
	}
}
