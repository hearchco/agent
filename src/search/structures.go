package search

type DomPaths struct {
	Result      string // div
	Link        string // a href
	Title       string // heading
	Description string // paragraph
	NextPage    string // button
}

type Options struct {
	UserAgent      string
	Limit          int
	ProxyAddr      string
	FollowNextPage bool
}

type Result struct {
	Rank        int
	URL         string
	Title       string
	Description string
}
