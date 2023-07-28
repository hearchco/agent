package search

const DefaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"

type DomQuery struct {
	Result      string
	Link        string
	Title       string
	Description string
	NextPage    string
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
