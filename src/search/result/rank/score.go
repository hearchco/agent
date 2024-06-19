package rank

import (
	"math"

	"github.com/hearchco/agent/src/config"
	"github.com/hearchco/agent/src/search/result"
)

// Calculates and sets scores for all results.
func (r Results) calculateScores(rconf config.CategoryRanking) {
	for _, res := range r {
		res.SetScore(calculateScore(res, rconf))
	}
}

// Calculates the score for one result.
func calculateScore(res result.Result, rconf config.CategoryRanking) float64 {
	var retRankScore float64 = 0
	for _, er := range res.EngineRanks() {
		eng := rconf.Engines[er.SearchEngine().ToLower()]
		retRankScore += (100.0/math.Pow(float64(er.Rank())*rconf.A+rconf.B, rconf.REXP)*rconf.C+rconf.D)*eng.Mul + eng.Const
	}

	retRankScore /= float64(len(res.EngineRanks()))
	timesReturnedScore := math.Log(float64(len(res.EngineRanks()))*rconf.TRA+rconf.TRB)*10*rconf.TRC + rconf.TRD
	score := retRankScore + timesReturnedScore

	return score
}
