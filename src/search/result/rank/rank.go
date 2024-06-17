package rank

import (
	"sort"

	"github.com/hearchco/agent/src/config"
	"github.com/hearchco/agent/src/search/result"
)

type Results []result.Result

// Cast []result.Result to Results, call rank() and return []result.Result slice.
func Rank(results []result.Result, rconf config.CategoryRanking) []result.Result {
	resType := make(Results, 0, len(results))
	for _, res := range results {
		resType = append(resType, res)
	}

	// Rank the results.
	resType.rank(rconf)

	rankedRes := make([]result.Result, 0, len(resType))
	for _, res := range resType {
		rankedRes = append(rankedRes, res)
	}
	return rankedRes
}

// Calculates the Score, sorts by it and then populates the Rank field of every result.
func (r Results) rank(rconf config.CategoryRanking) {
	// Fill Rank field for every EngineRank.
	r.fillEngineRankRank()

	// Calculate and set scores.
	r.calculateScores(rconf)

	// Sort slice by score.
	sort.Sort(ByScore(r))

	// Set correct ranks, by iterating over the sorted slice.
	r.correctRanks()
}

func (r Results) correctRanks() {
	for i, res := range r {
		res.SetRank(i + 1)
	}
}
