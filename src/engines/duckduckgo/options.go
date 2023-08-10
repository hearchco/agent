package duckduckgo

import "github.com/tminaorg/brzaguza/src/structures"

var Info structures.SEInfo = structures.SEInfo{
	Domain:     "lite.duckduckgo.com",
	Name:       "DuckDuckGo",
	URL:        "https://lite.duckduckgo.com/lite/",
	ResPerPage: 10,
	Crawlers:   []structures.EngineName{structures.Bing},
}

var dompaths structures.SEDOMPaths = structures.SEDOMPaths{
	Result:      "div.filters > table > tbody",
	Link:        "a.result-link",
	Title:       "td > a.result-link",
	Description: "td.result-snippet",
	NextPage:    "p.b_algoSlug",
}
