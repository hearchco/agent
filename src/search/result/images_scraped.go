package result

import (
	"github.com/hearchco/agent/src/utils/moreurls"
	"github.com/rs/zerolog/log"
)

type ImagesScraped struct {
	GeneralScraped

	originalSize  scrapedImageFormat
	thumbnailSize scrapedImageFormat
	thumbnailURL  string
	sourceName    string
	sourceURL     string
}

func (r ImagesScraped) OriginalSize() scrapedImageFormat {
	if r.originalSize.height == 0 || r.originalSize.width == 0 {
		log.Panic().
			Int("height", r.originalSize.height).
			Int("width", r.originalSize.width).
			Msg("OriginalSize is zero")
		// ^PANIC - Assert because the OriginalSize should never be zero.
	}

	return r.originalSize
}

func (r ImagesScraped) ThumbnailSize() scrapedImageFormat {
	if r.thumbnailSize.height == 0 || r.thumbnailSize.width == 0 {
		log.Panic().
			Int("height", r.thumbnailSize.height).
			Int("width", r.thumbnailSize.width).
			Msg("ThumbnailSize is zero")
		// ^PANIC - Assert because the ThumbnailSize should never be zero.
	}

	return r.thumbnailSize
}

func (r ImagesScraped) ThumbnailURL() string {
	if r.thumbnailURL == "" {
		log.Panic().Msg("ThumbnailURL is empty")
		// ^PANIC - Assert because the ThumbnailURL should never be empty.
	}

	return r.thumbnailURL
}

func (r ImagesScraped) SourceName() string {
	return r.sourceName
}

func (r ImagesScraped) SourceURL() string {
	return r.sourceURL
}

func (r ImagesScraped) Convert(erCap int) Result {
	engineRanks := make([]Rank, 0, erCap)
	engineRanks = append(engineRanks, r.Rank().Convert())
	return &Images{
		imagesJSON{
			General{
				generalJSON{
					URL:         r.URL(),
					FQDN:        moreurls.FQDN(r.URL()),
					Title:       r.Title(),
					Description: r.Description(),
					EngineRanks: engineRanks,
				},
			},
			r.OriginalSize().Convert(),
			r.ThumbnailSize().Convert(),
			r.ThumbnailURL(),
			r.SourceName(),
			r.SourceURL(),
		},
	}
}

type scrapedImageFormat struct {
	height int
	width  int
}

func (i scrapedImageFormat) GetHeight() int {
	if i.height == 0 {
		log.Panic().Msg("Height is zero")
		// ^PANIC - Assert because the Height should never be zero.
	}

	return i.height
}

func (i scrapedImageFormat) GetWidth() int {
	if i.width == 0 {
		log.Panic().Msg("Width is zero")
		// ^PANIC - Assert because the Width should never be zero.
	}

	return i.width
}

func (i scrapedImageFormat) Convert() ImageFormat {
	return ImageFormat{
		Height: i.height,
		Width:  i.width,
	}
}
