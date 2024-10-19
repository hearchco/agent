package result

import (
	"slices"
	"sync"

	"github.com/hearchco/agent/src/search/engines"
	"github.com/rs/zerolog/log"
)

type SuggestionConcMap struct {
	enabledEnginesLen int
	Mutex             sync.RWMutex
	Map               map[string]Suggestion
}

func NewSuggestionMap(enabledEnginesLen int) SuggestionConcMap {
	return SuggestionConcMap{
		enabledEnginesLen: enabledEnginesLen,
		Mutex:             sync.RWMutex{},
		Map:               make(map[string]Suggestion),
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
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	mapVal, exists := m.Map[val.Key()]
	if !exists {
		// Add the result to the map.
		m.Map[val.Key()] = val.Convert(m.enabledEnginesLen)
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
	m.Mutex.RLock()
	defer m.Mutex.RUnlock()

	suggestions := make([]Suggestion, 0, len(m.Map))
	responders := make([]engines.Name, 0, m.enabledEnginesLen)

	for _, sug := range m.Map {
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
