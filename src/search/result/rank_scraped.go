package result

import (
	"github.com/hearchco/agent/src/search/engines"
)

type RankScraped struct {
	searchEngine engines.Name
	rank         int
	page         int
	onPageRank   int
}

func (r RankScraped) SearchEngine() engines.Name {
	return r.searchEngine
}

func (r RankScraped) Rank() int {
	return r.rank
}

func (r RankScraped) Page() int {
	return r.page
}

func (r RankScraped) OnPageRank() int {
	return r.onPageRank
}

func (r RankScraped) Convert() Rank {
	return Rank{
		rankJSON: rankJSON{
			SearchEngine: r.searchEngine,
			Rank:         r.rank,
			Page:         r.page,
			OnPageRank:   r.onPageRank,
		},
	}
}

func NewRankScraped(searchEngine engines.Name, rank, page, onPageRank int) RankScraped {
	return RankScraped{
		searchEngine: searchEngine,
		rank:         rank,
		page:         page,
		onPageRank:   onPageRank,
	}
}
