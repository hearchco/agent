package bucket

import (
	"github.com/hearchco/hearchco/src/search/engines"
	"github.com/hearchco/hearchco/src/search/result"
)

// Returns the result made, and true if successful. If the result is not valid, false is returned.
func MakeSEResult(
	urll, title, desc string,
	seName engines.Name, sePage, seOnPageRank int,
) result.RetrievedResult {

	ser := result.RetrievedRank{
		SearchEngine: seName,
		Rank:         0,
		Page:         uint(sePage),
		OnPageRank:   uint(seOnPageRank),
	}

	res := result.RetrievedResult{
		URL:         urll,
		Title:       title,
		Description: desc,
		Rank:        ser,
	}

	return res
}

func MakeSEImageResult(
	urll, title, desc string,
	src, srcUrl, thmbUrl string,
	origH, origW, thmbH, thmbW int,
	seName engines.Name, sePage, seOnPageRank int,
) result.RetrievedResult {

	ser := result.RetrievedRank{
		SearchEngine: seName,
		Rank:         0,
		Page:         uint(sePage),
		OnPageRank:   uint(seOnPageRank),
	}

	res := result.RetrievedResult{
		URL:         urll,
		Title:       title,
		Description: desc,
		ImageResult: result.ImageResult{
			Original: result.ImageFormat{
				Height: uint(origH),
				Width:  uint(origW),
			},
			Thumbnail: result.ImageFormat{
				Height: uint(thmbH),
				Width:  uint(thmbW),
			},
			ThumbnailURL: thmbUrl,
			Source:       src,
			SourceURL:    srcUrl,
		},
		Rank: ser,
	}

	return res
}
