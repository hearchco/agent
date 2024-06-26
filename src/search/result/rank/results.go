package rank

import (
	"sort"

	"github.com/hearchco/agent/src/config"
	"github.com/hearchco/agent/src/search/result"
)

type Results []result.Result

// Calculates the Score, sorts by it and then populates the Rank field of every result.
func (r Results) Rank(rconf config.CategoryRanking) {
	// Fill Rank field for every EngineRank.
	r.fillEngineRankRank()

	// Calculate and set scores.
	r.calculateScores(rconf)

	// Sort slice by score.
	sort.Sort(ByScore[result.Result](r))

	// Set correct ranks, by iterating over the sorted slice.
	r.correctRanks()
}

func (r Results) correctRanks() {
	for i, res := range r {
		res.SetRank(i + 1)
	}
}
