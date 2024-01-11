package duckduckgo

import "github.com/hearchco/hearchco/src/engines"

var Info engines.Info = engines.Info{
	Domain:         "lite.duckduckgo.com",
	Name:           engines.DUCKDUCKGO,
	URL:            "https://lite.duckduckgo.com/lite/",
	ResultsPerPage: 10,
	Crawlers:       []engines.Name{engines.BING},
}

var dompaths engines.DOMPaths = engines.DOMPaths{
	ResultsContainer: "div.filters > table > tbody",
	Link:             "a.result-link",
	Title:            "td > a.result-link",
	Description:      "td.result-snippet",
}

var Support engines.SupportedSettings = engines.SupportedSettings{
	Locale: true,
}
