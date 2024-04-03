package result

import "github.com/hearchco/hearchco/src/search/engines"

// variables are 1-indexed
// Information about what Rank a result was on some Search Engine
type RetrievedRank struct {
	SearchEngine engines.Name `json:"search_engine"`
	Rank         uint         `json:"rank"`
	Page         uint         `json:"page"`
	OnPageRank   uint         `json:"on_page_rank"`
}

// The info a Search Engine returned about some Result
type RetrievedResult struct {
	URL         string        `json:"url"`
	URLHash     string        `json:"url_hash,omitempty"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	ImageResult ImageResult   `json:"image_result"`
	Rank        RetrievedRank `json:"rank"`
}
