package structures

type SEInfo struct {
	Domain     string
	Name       string
	URL        string
	ResPerPage int
}

type SEOptions struct {
	UserAgent     string
	MaxPages      int
	ProxyAddr     string
	JustFirstPage bool
	VisitPages    bool
}

type SEDOMPaths struct {
	ResultsContainer string
	Result           string // div
	Link             string // a href
	Title            string // heading
	Description      string // paragraph
	NextPage         string // button
}

type Engine string

const (
	Google     Engine = "google" // needs to be toLower
	Mojeek     Engine = "mojeek"
	DuckDuckGo Engine = "duckduckgo"
	Qwant      Engine = "qwant"
	Etools     Engine = "etools"
	Swisscows  Engine = "swisscows"
	Brave      Engine = "brave"
	Bing       Engine = "bing"
	Startpage  Engine = "startpage"
	Yandex     Engine = "yandex" // needed for crawler types
)
