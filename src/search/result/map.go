package result

import (
	"slices"
	"sync"

	"github.com/hearchco/agent/src/search/engines"
)

type ConcurrentMap struct {
	Mutex sync.RWMutex
	Map   map[string]Result
}

func Map() ConcurrentMap {
	return ConcurrentMap{
		Map: make(map[string]Result),
	}
}

func (r *ConcurrentMap) ExtractResultsAndResponders(enabledEnginesLen, titleLen, descLen int) ([]Result, []engines.Name) {
	r.Mutex.RLock()

	results := make([]Result, 0, len(r.Map))
	responders := make([]engines.Name, 0, enabledEnginesLen)

	for _, res := range r.Map {
		newRes := res.Shorten(titleLen, descLen)
		newRes.ShrinkEngineRanks()
		results = append(results, newRes)
		for _, rank := range res.EngineRanks() {
			if !slices.Contains(responders, rank.SearchEngine()) {
				responders = append(responders, rank.SearchEngine())
			}
		}
	}

	r.Mutex.RUnlock()

	return results, responders
}
