package search

import (
	"github.com/rs/zerolog/log"

	"github.com/hearchco/agent/src/search/engines"
	"github.com/hearchco/agent/src/search/result"
)

func createReceiver(engChan chan chan result.ResultScraped, concMap *result.ConcurrentMap, enabledEnginesLen int) {
	for resChan := range engChan {
		go func() {
			for recVal := range resChan {
				if recVal.Rank().SearchEngine().String() == "" || recVal.Rank().SearchEngine() == engines.UNDEFINED {
					log.Panic().
						Str("engine", recVal.Rank().SearchEngine().String()).
						Msg("Received a result with an undefined search engine")
					// ^PANIC - Assert because it should never happen.
				}

				// Lock the map due to modifications.
				concMap.Mutex.Lock()

				mapVal, exists := concMap.Map[recVal.Key()]
				if !exists {
					// Add the result to the map.
					concMap.Map[recVal.Key()] = recVal.Convert(enabledEnginesLen)
				} else {
					var alreadyIn *result.Rank

					// Check if the engine rank is already in the result.
					for i, er := range mapVal.EngineRanks() {
						if recVal.Rank().SearchEngine() == er.SearchEngine() {
							alreadyIn = &mapVal.EngineRanks()[i]
							break
						}
					}

					// Update the result if the new rank is better.
					if alreadyIn == nil {
						mapVal.AppendEngineRanks(recVal.Rank().Convert())
					} else {
						alreadyIn.UpgradeIfBetter(recVal.Rank().Convert())
					}

					// Update the description if the new description is longer.
					if len(mapVal.Description()) < len(recVal.Description()) {
						mapVal.SetDescription(recVal.Description())
					}
				}

				// Unlock the map.
				concMap.Mutex.Unlock()
			}
		}()
	}
}

func createSuggestionsReceiver(engChan chan chan result.SuggestionScraped, concMap *result.SuggestionConcurrentMap, enabledEnginesLen int) {
	for resChan := range engChan {
		go func() {
			for recVal := range resChan {
				if recVal.Rank().SearchEngine().String() == "" || recVal.Rank().SearchEngine() == engines.UNDEFINED {
					log.Panic().
						Str("engine", recVal.Rank().SearchEngine().String()).
						Msg("Received a result with an undefined search engine")
					// ^PANIC - Assert because it should never happen.
				}

				// Lock the map due to modifications.
				concMap.Mutex.Lock()

				mapVal, exists := concMap.Map[recVal.Key()]
				if !exists {
					// Add the result to the map.
					concMap.Map[recVal.Key()] = recVal.Convert(enabledEnginesLen)
				} else {
					var alreadyIn *result.RankSimple

					// Check if the engine rank is already in the result.
					for i, er := range mapVal.EngineRanks() {
						if recVal.Rank().SearchEngine() == er.SearchEngine() {
							alreadyIn = &mapVal.EngineRanks()[i]
							break
						}
					}

					// Update the result if the new rank is better.
					if alreadyIn == nil {
						mapVal.AppendEngineRanks(recVal.Rank().Convert())
					} else {
						alreadyIn.UpgradeIfBetter(recVal.Rank().Convert())
					}
				}

				// Unlock the map.
				concMap.Mutex.Unlock()
			}
		}()
	}
}
