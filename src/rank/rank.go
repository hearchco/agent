package rank

import (
	"sort"

	"github.com/tminaorg/brzaguza/src/bucket/result"
	"github.com/tminaorg/brzaguza/src/engines"
)

func SetRank(result *result.Result) {
	result.Rank = result.EngineRanks[0].Rank
}

func Rank(resMap map[string]*result.Result) []result.Result {
	results := make([]result.Result, 0, len(resMap))
	for _, res := range resMap {
		res.EngineRanks = res.EngineRanks[0:res.TimesReturned:res.TimesReturned]
		results = append(results, *res)
	}

	FillRetrievedRank(results)

	for ind := range results {
		SetRank(&(results[ind]))
	}

	sort.Sort(ByRank(results))

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
			results[el.ArrInd].EngineRanks[el.RRInd].Rank = rnk + 1
		}
	}
}
