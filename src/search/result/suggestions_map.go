package result

import (
	"slices"
	"sync"

	"github.com/hearchco/agent/src/search/engines"
	"github.com/rs/zerolog/log"
)

type SuggestionConcMap struct {
	enabledEnginesLen int
	mutex             sync.RWMutex
	mapp              map[string]Suggestion
}

func NewSuggestionMap(enabledEnginesLen int) SuggestionConcMap {
	return SuggestionConcMap{
		enabledEnginesLen: enabledEnginesLen,
		mutex:             sync.RWMutex{},
		mapp:              make(map[string]Suggestion),
	}
}

func (m *SuggestionConcMap) AddOrUpgrade(val SuggestionScraped) {
	if val.Rank().SearchEngine().String() == "" || val.Rank().SearchEngine() == engines.UNDEFINED {
		log.Panic().
			Str("engine", val.Rank().SearchEngine().String()).
			Msg("Received a suggestion with an undefined search engine")
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
		var alreadyIn *RankSimple

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
	}
}

func (m *SuggestionConcMap) ExtractWithResponders() ([]Suggestion, []engines.Name) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	suggestions := make([]Suggestion, 0, len(m.mapp))
	responders := make([]engines.Name, 0, m.enabledEnginesLen)

	for _, sug := range m.mapp {
		sug.ShrinkEngineRanks()
		suggestions = append(suggestions, sug)
		for _, rank := range sug.EngineRanks() {
			if !slices.Contains(responders, rank.SearchEngine()) {
				responders = append(responders, rank.SearchEngine())
			}
		}
	}

	return suggestions, responders
}
