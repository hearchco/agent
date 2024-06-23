package rank

import (
	"math"

	"github.com/hearchco/agent/src/config"
)

// Calculates and sets scores for all results.
func (r Results) calculateScores(rconf config.CategoryRanking) {
	for _, res := range r {
		res.SetScore(calculateScore(res, rconf))
	}
}

// Calculates and sets scores for all results.
func (s Suggestions) calculateScores(rconf config.CategoryRanking) {
	for i := range s {
		sug := &s[i]
		sug.SetScore(calculateScore(sug, rconf))
	}
}

// Calculates the score for one result.
func calculateScore[T ranker](val scoreEngineRanker[T], rconf config.CategoryRanking) float64 {
	var retRankScore float64 = 0
	for _, er := range val.EngineRanks() {
		eng := rconf.Engines[er.SearchEngine().String()]
		retRankScore += (100.0/math.Pow(float64(er.Rank())*rconf.A+rconf.B, rconf.REXP)*rconf.C+rconf.D)*eng.Mul + eng.Const
	}

	retRankScore /= float64(len(val.EngineRanks()))
	timesReturnedScore := math.Log(float64(len(val.EngineRanks()))*rconf.TRA+rconf.TRB)*10*rconf.TRC + rconf.TRD
	score := retRankScore + timesReturnedScore

	return score
}
