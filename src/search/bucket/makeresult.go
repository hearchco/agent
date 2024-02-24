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

func MakeSEImageResult(urll, title, desc, src, srcUrl string, orig, thmb result.Image, thmbUrl string, searchEngineName engines.Name, sePage, seOnPageRank int) *result.RetrievedResult {
	ser := result.RetrievedRank{
		SearchEngine: searchEngineName,
		Rank:         0,
		Page:         uint(sePage),
		OnPageRank:   uint(seOnPageRank),
	}
	res := result.RetrievedResult{
		URL:         urll,
		Title:       title,
		Description: desc,
		ImageResult: result.ImageResult{
			Original:     orig,
			Thumbnail:    thmb,
			ThumbnailURL: thmbUrl,
			Source:       src,
			SourceURL:    srcUrl,
		},
		Rank: ser,
	}
	return &res
}
