package rank

import (
	"sort"

	"github.com/hearchco/hearchco/src/config"
	"github.com/hearchco/hearchco/src/search/result"
)

func Rank(resMap map[string]*result.Result, rconf config.Ranking) []result.Result {
	results := make([]result.Result, 0, len(resMap))
	for _, res := range resMap {
		// set res.EngineRanks slice's capacity to it's length
		res.EngineRanks = res.EngineRanks[:len(res.EngineRanks):len(res.EngineRanks)]
		results = append(results, *res)
	}

	fillRetrievedRank(results)

	for ind := range results {
		results[ind].Score = getScore(&results[ind], &rconf)
	}
	sort.Sort(ByScore(results))
	for ind := range results {
		results[ind].Rank = uint(ind + 1)
	}

	return results
}
