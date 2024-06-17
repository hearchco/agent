package rank

import (
	"sort"

	"github.com/hearchco/agent/src/search/engines"
	"github.com/hearchco/agent/src/search/result"
)

// calculates Rank value of every EngineRank for each Search Engine individually by using Page and OnPageRank to sort
func (res Results) fillEngineRankRank() {
	seEngineRanks := make([][]*result.Rank, len(engines.NameValues()))

	for _, r := range res {
		for i := range r.EngineRanks() {
			er := &r.EngineRanks()[i]
			seEngineRanks[er.SearchEngine()] = append(seEngineRanks[er.SearchEngine()], er)
		}
	}

	for _, seer := range seEngineRanks {
		sort.Sort(ByPageAndOnPageRank(seer))
		for i, er := range seer {
			er.SetRank(i + 1)
		}
	}
}
