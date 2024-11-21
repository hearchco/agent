package result

import (
	"fmt"

	"github.com/hearchco/agent/src/search/engines"
)

func ConstructResult(seName engines.Name, urll string, title string, description string, page int, onPageRank int) (WebScraped, error) {
	if urll == "" {
		return WebScraped{}, fmt.Errorf("invalid URL: empty")
	}

	if title == "" {
		return WebScraped{}, fmt.Errorf("invalid title: empty")
	}

	if page <= 0 {
		return WebScraped{}, fmt.Errorf("invalid page: %d", page)
	}

	if onPageRank <= 0 {
		return WebScraped{}, fmt.Errorf("invalid onPageRank: %d", onPageRank)
	}

	return WebScraped{
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
	}, nil
}

func ConstructImagesResult(
	seName engines.Name, urll string, title string, description string, page int, onPageRank int,
	originalHeight int, originalWidth int, thumbnailHeight int, thumbnailWidth int,
	thumbnailUrl string, sourceName string, sourceUrl string,
) (ImagesScraped, error) {
	res, err := ConstructResult(seName, urll, title, description, page, onPageRank)
	if err != nil {
		return ImagesScraped{}, err
	}

	if originalHeight <= 0 {
		return ImagesScraped{}, fmt.Errorf("invalid originalHeight: %d", originalHeight)
	}

	if originalWidth <= 0 {
		return ImagesScraped{}, fmt.Errorf("invalid originalWidth: %d", originalWidth)
	}

	if thumbnailHeight <= 0 {
		return ImagesScraped{}, fmt.Errorf("invalid thumbnailHeight: %d", thumbnailHeight)
	}

	if thumbnailWidth <= 0 {
		return ImagesScraped{}, fmt.Errorf("invalid thumbnailWidth: %d", thumbnailWidth)
	}

	if thumbnailUrl == "" {
		return ImagesScraped{}, fmt.Errorf("invalid thumbnailUrl: empty")
	}

	if sourceUrl == "" {
		return ImagesScraped{}, fmt.Errorf("invalid sourceUrl: empty")
	}

	return ImagesScraped{
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
	}, nil
}
