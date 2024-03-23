package result

import (
	"github.com/gocolly/colly/v2"
)

// Everything about some Result, calculated and compiled from multiple search engines
// The URL is the primary key
type Result struct {
	URL         string
	URLHash     string // don't store in cache since private salt might change
	Rank        uint
	Score       float64
	Title       string
	Description string
	ImageResult ImageResult
	EngineRanks []RetrievedRank
	Response    *colly.Response // don't store in cache since it's too big
}

type ImageResult struct {
	Original         ImageFormat
	Thumbnail        ImageFormat
	ThumbnailURL     string
	ThumbnailURLHash string // don't store in cache since private salt might change
	Source           string
	SourceURL        string
}

type ImageFormat struct {
	Height uint
	Width  uint
}
