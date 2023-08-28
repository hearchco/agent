package engines

// variables are 1-indexed
// Information about what Rank a result was on some Search Engine
type RetrievedRank struct {
	SearchEngine string // this should be changed to Name
	Rank         int
	Page         int
	OnPageRank   int
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
	Name           string
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
}

type Name string

const (
	Google     Name = "google" // needs to be toLower
	Mojeek     Name = "mojeek"
	DuckDuckGo Name = "duckduckgo"
	Qwant      Name = "qwant"
	Etools     Name = "etools"
	Swisscows  Name = "swisscows"
	Brave      Name = "brave"
	Bing       Name = "bing"
	Startpage  Name = "startpage"
	Yandex     Name = "yandex" // needed for crawler types
	Yep        Name = "yep"
)
