package structures

type SupportedSettings struct {
	Locale                  bool
	SafeSearch              bool
	Mobile                  bool
	RequestedResultsPerPage bool
}

type SEInfo struct {
	Domain         string
	Name           string
	URL            string
	ResultsPerPage int
	Crawlers       []EngineName
}

type SEDOMPaths struct {
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
	Locale        string
	SafeSearch    bool
	Mobile        bool
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
