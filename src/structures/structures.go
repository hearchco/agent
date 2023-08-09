package structures

import (
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
)

// Delegates Timeout, PageTimeout to colly.Collector.SetRequestTimeout(); Note: See https://github.com/gocolly/colly/issues/644
// Delegates Delay, RandomDelay, Parallelism to colly.Collector.Limit()
type Timings struct {
	Timeout     time.Duration
	PageTimeout time.Duration
	Delay       time.Duration
	RandomDelay time.Duration
	Parallelism int
}

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
