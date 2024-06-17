package result

type ImagesScraped struct {
	GeneralScraped

	originalSize  scrapedImageFormat
	thumbnailSize scrapedImageFormat
	thumbnailURL  string
	sourceName    string
	sourceURL     string
}

func (r ImagesScraped) OriginalSize() scrapedImageFormat {
	return r.originalSize
}

func (r ImagesScraped) ThumbnailSize() scrapedImageFormat {
	return r.thumbnailSize
}

func (r ImagesScraped) ThumbnailURL() string {
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
	return i.height
}

func (i scrapedImageFormat) GetWidth() int {
	return i.width
}

func (i scrapedImageFormat) Convert() ImageFormat {
	return ImageFormat{
		Height: i.height,
		Width:  i.width,
	}
}
