package result

import (
	"github.com/gocolly/colly/v2"
)

type ImageFormat struct {
	Height uint `json:"height"`
	Width  uint `json:"width"`
}

type ImageResult struct {
	Original         ImageFormat `json:"original"`
	Thumbnail        ImageFormat `json:"thumbnail"`
	ThumbnailURL     string      `json:"thumbnail_url"`
	ThumbnailURLHash string      `json:"thumbnail_url_hash,omitempty"`
	Source           string      `json:"source"`
	SourceURL        string      `json:"source_url"`
}

// Everything about some Result, calculated and compiled from multiple search engines
// The URL is the primary key
type Result struct {
	URL         string          `json:"url"`
	URLHash     string          `json:"url_hash,omitempty"`
	Rank        uint            `json:"rank"`
	Score       float64         `json:"score"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	EngineRanks []RetrievedRank `json:"engine_ranks"`
	ImageResult ImageResult     `json:"image_result"`
	Response    *colly.Response `json:"-"`
}
