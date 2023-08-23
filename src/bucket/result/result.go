package result

import (
	"github.com/gocolly/colly/v2"
	"github.com/tminaorg/brzaguza/src/engines"
)

// Everything about some Result, calculated and compiled from multiple search engines
// The URL is the primary key
type Result struct {
	URL           string
	Rank          int
	Title         string
	Description   string
	EngineRanks   []engines.RetrievedRank
	TimesReturned int
	Response      *colly.Response
}
