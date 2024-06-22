package result

import (
	"slices"
	"sync"

	"github.com/hearchco/agent/src/search/engines"
)

type SuggestionConcurrentMap struct {
	enabledEnginesLen int
	Mutex             sync.RWMutex
	Map               map[string]Suggestion
}

func SuggestionMap(enabledEnginesLen int) SuggestionConcurrentMap {
	return SuggestionConcurrentMap{
		enabledEnginesLen: enabledEnginesLen,
		Mutex:             sync.RWMutex{},
		Map:               make(map[string]Suggestion),
	}
}

func (r *SuggestionConcurrentMap) ExtractSuggestionsAndResponders() ([]Suggestion, []engines.Name) {
	r.Mutex.RLock()

	suggestions := make([]Suggestion, 0, len(r.Map))
	responders := make([]engines.Name, 0, r.enabledEnginesLen)

	for _, sug := range r.Map {
		sug.ShrinkEngineRanks()
		suggestions = append(suggestions, sug)
		for _, rank := range sug.EngineRanks() {
			if !slices.Contains(responders, rank.SearchEngine()) {
				responders = append(responders, rank.SearchEngine())
			}
		}
	}

	r.Mutex.RUnlock()

	return suggestions, responders
}
