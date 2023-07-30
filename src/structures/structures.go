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

type Result struct {
	URL         string
	Rank        int
	SEPageRank  int
	SEPage      int
	Title       string
	Description string
	Response    *colly.Response
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
