package bucket

import (
	"github.com/hearchco/hearchco/src/search/engines"
	"github.com/hearchco/hearchco/src/search/result"
)

func MakeSEResult(urll string, title string, description string, searchEngineName engines.Name, sePage int, seOnPageRank int) *result.RetrievedResult {
	ser := result.RetrievedRank{
		SearchEngine: searchEngineName,
		Rank:         0,
		Page:         uint(sePage),
		OnPageRank:   uint(seOnPageRank),
	}
	res := result.RetrievedResult{
		URL:         urll,
		Title:       title,
		Description: description,
		Rank:        ser,
	}
	return &res
}
