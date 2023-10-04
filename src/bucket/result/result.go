package result

import (
	"encoding/json"

	"github.com/gocolly/colly/v2"
	"github.com/tminaorg/brzaguza/src/engines"
)

// Everything about some Result, calculated and compiled from multiple search engines
// The URL is the primary key
type Result struct {
	URL           string                  `json:"url"`
	Rank          uint                    `json:"rank"`
	Score         float64                 `json:"score"`
	Title         string                  `json:"title"`
	Description   string                  `json:"description"`
	EngineRanks   []engines.RetrievedRank `json:"engineRanks"`
	TimesReturned uint8                   `json:"timesReturned"`
	Response      *colly.Response         `json:"response"`
}

// MarshalJSON implements the json.Marshaler interface for Results
func (r *Result) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}

type Results []Result

// MarshalJSON implements the json.Marshaler interface for Results
func (r Results) MarshalJSON() ([]byte, error) {
	return json.Marshal([]Result(r))
}
