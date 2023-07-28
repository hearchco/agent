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
	Rank        int    `json:"rank"`
	URL         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (r Result) Hash() string {
	return r.URL
}

type ByRank []Result

func (r ByRank) Len() int           { return len(r) }
func (r ByRank) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r ByRank) Less(i, j int) bool { return r[i].Rank < r[j].Rank }
