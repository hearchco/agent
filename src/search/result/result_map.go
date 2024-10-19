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
	Mutex             sync.RWMutex
	Map               map[string]Result
}

func NewResultMap(enabledEnginesLen, titleLen, descLen int) ResultConcMap {
	return ResultConcMap{
		enabledEnginesLen: enabledEnginesLen,
		titleLen:          titleLen,
		descLen:           descLen,
		Mutex:             sync.RWMutex{},
		Map:               make(map[string]Result),
	}
}

func (m *ResultConcMap) AddOrUpgrade(val ResultScraped) {
	if val.Rank().SearchEngine().String() == "" || val.Rank().SearchEngine() == engines.UNDEFINED {
		log.Panic().
			Str("engine", val.Rank().SearchEngine().String()).
			Msg("Received a result with an undefined search engine")
		// ^PANIC - Assert because it should never happen.
	}

	// Lock the map due to modifications.
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	mapVal, exists := m.Map[val.Key()]
	if !exists {
		// Add the result to the map.
		m.Map[val.Key()] = val.Convert(m.enabledEnginesLen)
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

func (m *ResultConcMap) ExtractWithResponders() ([]Result, []engines.Name) {
	m.Mutex.RLock()
	defer m.Mutex.RUnlock()

	results := make([]Result, 0, len(m.Map))
	responders := make([]engines.Name, 0, m.enabledEnginesLen)

	for _, res := range m.Map {
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
