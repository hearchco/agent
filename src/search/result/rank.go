package result

import (
	"github.com/hearchco/agent/src/search/engines"
)

type Rank struct {
	RankSimple

	rankJSON
}

type rankJSON struct {
	Page       int `json:"page"`
	OnPageRank int `json:"on_page_rank"`
}

func (r Rank) Page() int {
	return r.rankJSON.Page
}

func (r *Rank) SetPage(page, onPageRank int) {
	r.rankJSON.Page = page
	r.rankJSON.OnPageRank = onPageRank
}

func (r Rank) OnPageRank() int {
	return r.rankJSON.OnPageRank
}

func (r *Rank) SetOnPageRank(onPageRank int) {
	r.rankJSON.OnPageRank = onPageRank
}

func (r *Rank) UpgradeIfBetter(newR Rank) {
	if r.Page() > newR.Page() {
		r.SetPage(newR.Page(), newR.OnPageRank())
	} else if r.Page() == newR.Page() && r.OnPageRank() > newR.OnPageRank() {
		r.SetOnPageRank(newR.OnPageRank())
	}
}

func NewRank(searchEngine engines.Name, rank, page, onPageRank int) Rank {
	return Rank{
		RankSimple{
			rankSimpleJSON{
				SearchEngine: searchEngine,
				Rank:         rank,
			},
		},
		rankJSON{
			page,
			onPageRank,
		},
	}
}
