package engines

import "github.com/tminaorg/brzaguza/src/category"

// variables are 1-indexed
// Information about what Rank a result was on some Search Engine
type RetrievedRank struct {
	SearchEngine Name
	Rank         uint
	Page         uint
	OnPageRank   uint
}

// The info a Search Engine returned about some Result
type RetrievedResult struct {
	URL         string
	Title       string
	Description string
	Rank        RetrievedRank
}

type SupportedSettings struct {
	Locale                  bool
	SafeSearch              bool
	Mobile                  bool
	RequestedResultsPerPage bool
}

type Info struct {
	Domain         string
	Name           Name
	URL            string
	ResultsPerPage int
	Crawlers       []Name
}

type DOMPaths struct {
	ResultsContainer string
	Result           string // div
	Link             string // a href
	Title            string // heading
	Description      string // paragraph
}

type Options struct {
	UserAgent     string
	MaxPages      int
	ProxyAddr     string
	JustFirstPage bool
	VisitPages    bool
	Locale        string //format: en-US
	SafeSearch    bool
	Mobile        bool
	Category      category.Name
}
