package duckduckgo

import "github.com/tminaorg/brzaguza/src/engines"

var Info engines.Info = engines.Info{
	Domain:         "lite.duckduckgo.com",
	Name:           "DuckDuckGo",
	URL:            "https://lite.duckduckgo.com/lite/",
	ResultsPerPage: 10,
	Crawlers:       []engines.Name{engines.Bing},
}

var dompaths engines.DOMPaths = engines.DOMPaths{
	ResultsContainer: "div.filters > table > tbody",
	Link:             "a.result-link",
	Title:            "td > a.result-link",
	Description:      "td.result-snippet",
}

var Support engines.SupportedSettings = engines.SupportedSettings{}
