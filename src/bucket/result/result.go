package result

import (
	"encoding/json"
	"fmt"

	"github.com/gocolly/colly/v2"
	"github.com/tminaorg/brzaguza/src/engines"
)

// Everything about some Result, calculated and compiled from multiple search engines
// The URL is the primary key
type Result struct {
	URL           string                  `json:"url"`
	Rank          int                     `json:"rank"`
	Title         string                  `json:"title"`
	Description   string                  `json:"description"`
	EngineRanks   []engines.RetrievedRank `json:"engineRanks"`
	TimesReturned int                     `json:"timesReturned"`
	Response      *colly.Response         `json:"response"`
}

// MarshalJSON implements the json.Marshaler interface for Results
func (r Result) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf("%v", r))
}

// UnmarshalJSON implements the json.Unmarshaler interface for Results
func (r *Result) UnmarshalJSON(data []byte) error {
	var s Result
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("data should be of type Result, got %s", data)
	}

	*r = s
	return nil
}

type Results []Result

// MarshalJSON implements the json.Marshaler interface for Results
func (r Results) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf("%v", r))
}

// UnmarshalJSON implements the json.Unmarshaler interface for Results
func (r *Results) UnmarshalJSON(data []byte) error {
	var s Results
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("data should be of type Results, got %s", data)
	}

	*r = s
	return nil
}
