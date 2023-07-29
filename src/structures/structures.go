package structures

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
	Limit         int
	ProxyAddr     string
	JustFirstPage bool
}

type Result struct {
	Rank        int
	URL         string
	Title       string
	Description string
}

type ResultRank struct {
	URL  string
	Rank int
}
