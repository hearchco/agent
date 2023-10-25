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

func firstN(str string, n int) string {
	v := []rune(str)
	if n >= len(v) {
		return str
	}
	return string(v[:n])
}

func Shorten(results []Result) []Result {
	resultsShort := make([]Result, len(results))
	copy(resultsShort, results)
	for _, r := range resultsShort {
		r.Description = firstN(r.Description, 400)
	}
	return resultsShort
}
