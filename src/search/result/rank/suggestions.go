package rank

import (
	"sort"

	"github.com/hearchco/agent/src/search/category"
	"github.com/hearchco/agent/src/search/result"
)

type Suggestions []result.Suggestion

// Calculates the Score, sorts by it and then populates the Rank field of every result.
func (s Suggestions) Rank(rconf category.Ranking) {
	// Calculate and set scores.
	s.calculateScores(rconf)

	// Sort slice by score.
	sort.Sort(ByScore[result.Suggestion](s))

	// Set correct ranks, by iterating over the sorted slice.
	s.correctRanks()
}

func (s Suggestions) correctRanks() {
	for i := range s {
		sug := &s[i]
		sug.SetRank(i + 1)
	}
}
