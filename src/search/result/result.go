package result

import (
	"github.com/gocolly/colly/v2"
)

type ImageFormat struct {
	Height uint
	Width  uint
}

type ImageResult struct {
	Original     ImageFormat
	Thumbnail    ImageFormat
	ThumbnailURL string
	Source       string
	SourceURL    string
}

// Everything about some Result, calculated and compiled from multiple search engines
// The URL is the primary key
type Result struct {
	URL           string
	Rank          uint
	Score         float64
	Title         string
	Description   string
	EngineRanks   []RetrievedRank
	TimesReturned uint8
	ImageResult   ImageResult
	Response      *colly.Response
}
