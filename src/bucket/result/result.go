package result

import (
	"github.com/gocolly/colly/v2"
	"github.com/tminaorg/brzaguza/src/engines"
)

// Everything about some Result, calculated and compiled from multiple search engines
// The URL is the primary key
type Result struct {
	URL           string
	Rank          uint
	Score         float64
	Title         string
	Description   string
	EngineRanks   []engines.RetrievedRank
	TimesReturned uint8
	Response      *colly.Response
}
