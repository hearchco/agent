package rank

import (
	"sort"

	"github.com/hearchco/hearchco/src/search/engines"
	"github.com/hearchco/hearchco/src/search/result"
)

type RankFiller struct {
	ArrInd  int
	RetRank result.RetrievedRank
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
