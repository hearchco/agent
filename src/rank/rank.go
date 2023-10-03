package rank

import (
	"math"
	"sort"

	"github.com/tminaorg/brzaguza/src/bucket/result"
	"github.com/tminaorg/brzaguza/src/config"
	"github.com/tminaorg/brzaguza/src/engines"
)

func GetScore(result *result.Result, rconf *config.Ranking) float64 {
	retRankScore := float64(0)
	for _, er := range result.EngineRanks {
		seMul := rconf.Engines[er.SearchEngine.ToLower()].Mul
		seConst := rconf.Engines[er.SearchEngine.ToLower()].Const //these 2 could be preproced into array
		retRankScore += (100.0/math.Pow(float64(er.Rank)*rconf.A+rconf.B, rconf.REXP)*rconf.C+rconf.D)*seMul + seConst
	}
	retRankScore /= float64(result.TimesReturned)

	timesReturnedScore := math.Log(float64(result.TimesReturned)*rconf.TRA+rconf.TRB)*10*rconf.TRC + rconf.TRD

	score := retRankScore + timesReturnedScore
	return score
}

func Rank(resMap map[string]*result.Result, rconf *config.Ranking) []result.Result {
	results := make([]result.Result, 0, len(resMap))
	for _, res := range resMap {
		res.EngineRanks = res.EngineRanks[0:res.TimesReturned:res.TimesReturned]
		results = append(results, *res)
	}

	FillRetrievedRank(results)

	for ind := range results {
		results[ind].Score = GetScore(&(results[ind]), rconf)
	}
	sort.Sort(ByScore(results))
	for ind := range results {
		results[ind].Rank = uint(ind + 1)
	}

	return results
}

type RankFiller struct {
	ArrInd  int
	RetRank engines.RetrievedRank
	RRInd   int
}

func FillRetrievedRank(results []result.Result) {
	engResults := make([][]RankFiller, len(engines.NameValues()))
	for arrind, res := range results {
		for rrind, er := range res.EngineRanks {
			rf := RankFiller{
				ArrInd:  arrind,
				RetRank: er,
				RRInd:   rrind,
			}
			engResults[er.SearchEngine] = append(engResults[er.SearchEngine], rf)
		}
	}

	for _, engRes := range engResults {
		sort.Sort(ByRetrievedRank(engRes))

		for rnk, el := range engRes {
			results[el.ArrInd].EngineRanks[el.RRInd].Rank = uint(rnk + 1)
		}
	}
}
