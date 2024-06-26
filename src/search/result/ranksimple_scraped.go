package result

import (
	"github.com/hearchco/agent/src/search/engines"
)

type RankSimpleScraped struct {
	searchEngine engines.Name
	rank         int
}

func (r RankSimpleScraped) SearchEngine() engines.Name {
	return r.searchEngine
}

func (r RankSimpleScraped) Rank() int {
	return r.rank
}

func (r RankSimpleScraped) Convert() RankSimple {
	return RankSimple{
		rankSimpleJSON{
			r.searchEngine,
			r.rank,
		},
	}
}

func NewRankSimpleScraped(searchEngine engines.Name, rank int) RankSimpleScraped {
	return RankSimpleScraped{
		searchEngine,
		rank,
	}
}
