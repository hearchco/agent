package rank

import (
	"math"

	"github.com/hearchco/agent/src/search/category"
)

// Calculates and sets scores for all results.
func (r Results) calculateScores(rconf category.Ranking) {
	for _, res := range r {
		res.SetScore(calculateScore(res, rconf))
	}
}

// Calculates and sets scores for all results.
func (s Suggestions) calculateScores(rconf category.Ranking) {
	for i := range s {
		sug := &s[i]
		sug.SetScore(calculateScore(sug, rconf))
	}
}

// Calculates the score for one result.
func calculateScore[T ranker](val scoreEngineRanker[T], rconf category.Ranking) float64 {
	var rankScoreSum float64 = 0

	// Calculate the sum of the rank scores of all engines.
	// The rank score is dividing 100 to invert the priority (the lower the rank, the higher the score).
	for _, er := range val.EngineRanks() {
		eng := rconf.Engines[er.SearchEngine()]
		rankScoreSum += (100.0/math.Pow(float64(er.Rank())*rconf.RankMul+rconf.RankAdd, rconf.RankExp)*rconf.RankScoreMul+rconf.RankScoreAdd)*eng.Mul + eng.Add
	}

	// Calculate the average rank score from the sum.
	rankScoreAvg := rankScoreSum / float64(len(val.EngineRanks()))

	// Calculate a second score based on the number of times the result was returned.
	// Log is used to make the score less sensitive to the number of times returned.
	timesReturnedScore := math.Log(float64(len(val.EngineRanks()))*rconf.TimesReturnedMul+rconf.TimesReturnedAdd)*100*rconf.TimesReturnedScoreMul + rconf.TimesReturnedScoreAdd

	return rankScoreAvg + timesReturnedScore
}
