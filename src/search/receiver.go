package search

import (
	"sync"

	"github.com/hearchco/agent/src/search/result"
)

func createReceiver[T any](wg *sync.WaitGroup, valChan chan T, concMap result.ConcMapper[T]) {
	// Signal that the receiver is done.
	defer wg.Done()

	for recVal := range valChan {
		concMap.AddOrUpgrade(recVal)
	}
}
