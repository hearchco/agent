package result

import (
	"github.com/hearchco/agent/src/search/engines"
)

type RankScraped struct {
	RankSimpleScraped

	page       int
	onPageRank int
}

func (r RankScraped) Page() int {
	return r.page
}

func (r RankScraped) OnPageRank() int {
	return r.onPageRank
}

func (r RankScraped) Convert() Rank {
	rankSimple := r.RankSimpleScraped.Convert()
	return Rank{
		rankSimple,
		rankJSON{
			r.page,
			r.onPageRank,
		},
	}
}

func NewRankScraped(searchEngine engines.Name, rank, page, onPageRank int) RankScraped {
	rankSimpleScraped := NewRankSimpleScraped(searchEngine, rank)
	return RankScraped{
		rankSimpleScraped,
		page,
		onPageRank,
	}
}
