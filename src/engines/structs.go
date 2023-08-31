package engines

// variables are 1-indexed
// Information about what Rank a result was on some Search Engine
type RetrievedRank struct {
	SearchEngine Name
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
}

type Name string

const (
	Google     Name = "Google" // needs to be toLower
	Mojeek     Name = "Mojeek"
	DuckDuckGo Name = "DuckDuckGo"
	Qwant      Name = "Qwant"
	Etools     Name = "Etools"
	Swisscows  Name = "Swisscows"
	Brave      Name = "Brave"
	Bing       Name = "Bing"
	Startpage  Name = "Startpage"
	Yandex     Name = "Yandex" // needed for crawler types
	Yep        Name = "Yep"
	Presearch  Name = "Presearch"
)
