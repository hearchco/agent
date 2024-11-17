package result

import (
	"fmt"

	"github.com/hearchco/agent/src/search/engines"
)

func ConstructResult(seName engines.Name, urll string, title string, description string, page int, onPageRank int) (WebScraped, error) {
	res := WebScraped{
		url:         urll,
		title:       title,
		description: description,
		rank: RankScraped{
			RankSimpleScraped{
				searchEngine: seName,
				rank:         0, // This gets calculated when ranking the results.
			},
			page,
			onPageRank,
		},
	}

	if urll == "" {
		return res, fmt.Errorf("invalid URL: empty")
	}

	if title == "" {
		return res, fmt.Errorf("invalid title: empty")
	}

	if page <= 0 {
		return res, fmt.Errorf("invalid page: %d", page)
	}

	if onPageRank <= 0 {
		return res, fmt.Errorf("invalid onPageRank: %d", onPageRank)
	}

	return res, nil
}

func ConstructImagesResult(
	seName engines.Name, urll string, title string, description string, page int, onPageRank int,
	originalHeight int, originalWidth int, thumbnailHeight int, thumbnailWidth int,
	thumbnailUrl string, sourceName string, sourceUrl string,
) (ImagesScraped, error) {
	res, err := ConstructResult(seName, urll, title, description, page, onPageRank)
	imgres := ImagesScraped{
		WebScraped: res,

		originalSize: scrapedImageFormat{
			height: originalHeight,
			width:  originalWidth,
		},
		thumbnailSize: scrapedImageFormat{
			height: thumbnailHeight,
			width:  thumbnailWidth,
		},
		thumbnailURL: thumbnailUrl,
		sourceName:   sourceName,
		sourceURL:    sourceUrl,
	}
	if err != nil {
		return imgres, err
	}

	if originalHeight <= 0 {
		return imgres, fmt.Errorf("invalid originalHeight: %d", originalHeight)
	}

	if originalWidth <= 0 {
		return imgres, fmt.Errorf("invalid originalWidth: %d", originalWidth)
	}

	if thumbnailHeight <= 0 {
		return imgres, fmt.Errorf("invalid thumbnailHeight: %d", thumbnailHeight)
	}

	if thumbnailWidth <= 0 {
		return imgres, fmt.Errorf("invalid thumbnailWidth: %d", thumbnailWidth)
	}

	if thumbnailUrl == "" {
		return imgres, fmt.Errorf("invalid thumbnailUrl: empty")
	}

	if sourceUrl == "" {
		return imgres, fmt.Errorf("invalid sourceUrl: empty")
	}

	return imgres, nil
}
