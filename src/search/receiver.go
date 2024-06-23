package search

import (
	"sync"

	"github.com/rs/zerolog/log"

	"github.com/hearchco/agent/src/search/engines"
	"github.com/hearchco/agent/src/search/result"
)

func createReceiver[T any, R result.Ranker, V result.ConcReceiver[R]](wg *sync.WaitGroup, valChan chan V, concMap result.ConcMapper[T, V]) {
	// Signal that the receiver is done.
	defer wg.Done()

	for recVal := range valChan {
		if recVal.Rank().SearchEngine().String() == "" || recVal.Rank().SearchEngine() == engines.UNDEFINED {
			log.Panic().
				Str("engine", recVal.Rank().SearchEngine().String()).
				Msg("Received a result with an undefined search engine")
			// ^PANIC - Assert because it should never happen.
		}

		concMap.AddOrUpgrade(recVal)
	}
}
