package structures

type Relay struct {
	ResultChannel     chan Result
	RankChannel       chan ResultRank
	EngineDoneChannel chan bool
	ResultMap         map[string]*Result
}

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
	MaxPages      int
	ProxyAddr     string
	JustFirstPage bool
	VisitPages    bool
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

func (r Result) Hash() string {
	return r.URL
}

type ByRank []Result

func (r ByRank) Len() int           { return len(r) }
func (r ByRank) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r ByRank) Less(i, j int) bool { return r[i].Rank < r[j].Rank }
