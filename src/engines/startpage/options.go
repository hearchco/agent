package startpage

import "github.com/tminaorg/brzaguza/src/engines"

var Info engines.Info = engines.Info{
	Domain:         "www.startpage.com",
	Name:           engines.Startpage,
	URL:            "https://www.startpage.com/sp/search?q=",
	ResultsPerPage: 10,
	Crawlers:       []engines.Name{engines.Google},
}

var dompaths engines.DOMPaths = engines.DOMPaths{
	Result:      "section > div.w-gl__result",
	Link:        "a.result-link",
	Title:       "a.w-gl__result-title",
	Description: "p.w-gl__description",
}

var Support engines.SupportedSettings = engines.SupportedSettings{
	SafeSearch: true,
}
