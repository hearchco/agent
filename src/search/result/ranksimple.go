package result

import (
	"github.com/hearchco/agent/src/search/engines"
)

type RankSimple struct {
	rankSimpleJSON
}

type rankSimpleJSON struct {
	SearchEngine engines.Name `json:"search_engine"`
	Rank         int          `json:"rank"`
}

func (r RankSimple) SearchEngine() engines.Name {
	return r.rankSimpleJSON.SearchEngine
}

func (r RankSimple) Rank() int {
	return r.rankSimpleJSON.Rank
}

func (r *RankSimple) SetRank(rank int) {
	r.rankSimpleJSON.Rank = rank
}

func (r *RankSimple) UpgradeIfBetter(newR RankSimple) {
	if r.Rank() > newR.Rank() {
		r.SetRank(newR.Rank())
	}
}
