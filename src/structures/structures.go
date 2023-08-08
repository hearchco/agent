package structures

import (
	"sync"

	"github.com/gocolly/colly/v2"
)

type Relay struct {
	ResultMap         map[string]*Result
	Mutex             sync.RWMutex
	EngineDoneChannel chan bool
}

type DOMPaths struct {
	ResultsContainer string
	Result           string // div
	Link             string // a href
	Title            string // heading
	Description      string // paragraph
	NextPage         string // button
}

type Options struct {
	UserAgent     string
	MaxPages      int
	ProxyAddr     string
	JustFirstPage bool
	VisitPages    bool
}

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

/*
func (r Result) Hash() string {
	return r.URL
}
*/

type ByRank []Result

func (r ByRank) Len() int           { return len(r) }
func (r ByRank) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r ByRank) Less(i, j int) bool { return r[i].Rank < r[j].Rank }

type Engine string

const (
	Google     Engine = "google"
	Mojeek     Engine = "mojeek"
	DuckDuckGo Engine = "duckduckgo"
	Qwant      Engine = "qwant"
	Etools     Engine = "etools"
	Swisscows  Engine = "swisscows"
	Brave      Engine = "brave"
	Bing       Engine = "bing"
	Startpage  Engine = "startpage"
)
