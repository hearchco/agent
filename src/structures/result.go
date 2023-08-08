package structures

import (
	"github.com/gocolly/colly/v2"
)

// variables are 1-indexed
// Information about what Rank a result was on some Search Engine
type SERank struct {
	SearchEngine string
	Rank         int
	Page         int
	OnPageRank   int
}

// The info a Search Engine returned about some Result
type SEResult struct {
	URL         string
	Title       string
	Description string
	Rank        SERank
}

// Everything about some Result, calculated and compiled from multiple search engines
// The URL is the primary key
type Result struct {
	URL           string
	Rank          int
	Title         string
	Description   string
	SearchEngines []SERank
	TimesReturned int
	Response      *colly.Response
}

type ByRank []Result

func (r ByRank) Len() int           { return len(r) }
func (r ByRank) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r ByRank) Less(i, j int) bool { return r[i].Rank < r[j].Rank }
