package structures

type SEInfo struct {
	Domain     string
	Name       string
	URL        string
	ResPerPage int
	Crawlers   []EngineName
}

type SEDOMPaths struct {
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

type EngineName string

const (
	Google     EngineName = "google" // needs to be toLower
	Mojeek     EngineName = "mojeek"
	DuckDuckGo EngineName = "duckduckgo"
	Qwant      EngineName = "qwant"
	Etools     EngineName = "etools"
	Swisscows  EngineName = "swisscows"
	Brave      EngineName = "brave"
	Bing       EngineName = "bing"
	Startpage  EngineName = "startpage"
	Yandex     EngineName = "yandex" // needed for crawler types
)
