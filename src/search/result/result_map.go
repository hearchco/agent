package result

import (
	"slices"
	"sync"

	"github.com/hearchco/agent/src/search/engines"
	"github.com/rs/zerolog/log"
)

type ResultConcMap struct {
	enabledEnginesLen int
	titleLen, descLen int
	mutex             sync.RWMutex
	mapp              map[string]Result
}

func NewResultMap(enabledEnginesLen, titleLen, descLen int) ResultConcMap {
	return ResultConcMap{
		enabledEnginesLen: enabledEnginesLen,
		titleLen:          titleLen,
		descLen:           descLen,
		mutex:             sync.RWMutex{},
		mapp:              make(map[string]Result),
	}
}

// Passed as pointer because of the mutex.
func (m *ResultConcMap) AddOrUpgrade(val ResultScraped) {
	if val.Rank().SearchEngine().String() == "" || val.Rank().SearchEngine() == engines.UNDEFINED {
		log.Panic().
			Str("engine", val.Rank().SearchEngine().String()).
			Msg("Received a result with an undefined search engine")
		// ^PANIC - Assert because it should never happen.
	}

	// Lock the map due to modifications.
	m.mutex.Lock()
	defer m.mutex.Unlock()

	mapVal, exists := m.mapp[val.Key()]
	if !exists {
		// Add the result to the map.
		m.mapp[val.Key()] = val.Convert(m.enabledEnginesLen)
	} else {
		var alreadyIn *Rank

		// Check if the engine rank is already in the result.
		for i, er := range mapVal.EngineRanks() {
			if val.Rank().SearchEngine() == er.SearchEngine() {
				alreadyIn = &mapVal.EngineRanks()[i]
				break
			}
		}

		// Update the result if the new rank is better.
		if alreadyIn == nil {
			mapVal.AppendEngineRanks(val.Rank().Convert())
		} else {
			alreadyIn.UpgradeIfBetter(val.Rank().Convert())
		}

		// Update the description if the new description is longer.
		if len(mapVal.Description()) < len(val.Description()) {
			mapVal.SetDescription(val.Description())
		}
	}
}

// Passed as pointer because of the mutex.
func (m *ResultConcMap) ExtractWithResponders() ([]Result, []engines.Name) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	results := make([]Result, 0, len(m.mapp))
	responders := make([]engines.Name, 0, m.enabledEnginesLen)

	for _, res := range m.mapp {
		newRes := res.Shorten(m.titleLen, m.descLen)
		newRes.ShrinkEngineRanks()
		results = append(results, newRes)
		for _, rank := range res.EngineRanks() {
			if !slices.Contains(responders, rank.SearchEngine()) {
				responders = append(responders, rank.SearchEngine())
			}
		}
	}

	return results, responders
}
