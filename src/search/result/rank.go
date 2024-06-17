package result

import (
	"github.com/hearchco/agent/src/search/engines"
)

type Rank struct {
	rankJSON
}

type rankJSON struct {
	SearchEngine engines.Name `json:"search_engine"`
	Rank         int          `json:"rank"`
	Page         int          `json:"page"`
	OnPageRank   int          `json:"on_page_rank"`
}

func (r Rank) SearchEngine() engines.Name {
	return r.rankJSON.SearchEngine
}

func (r Rank) Rank() int {
	return r.rankJSON.Rank
}

func (r *Rank) SetRank(rank int) {
	r.rankJSON.Rank = rank
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

func NewRank(searchEngine engines.Name, rank, page, onPageRank int) Rank {
	return Rank{
		rankJSON{
			SearchEngine: searchEngine,
			Rank:         rank,
			Page:         page,
			OnPageRank:   onPageRank,
		},
	}
}
