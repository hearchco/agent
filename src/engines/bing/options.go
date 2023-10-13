package bing

import "github.com/tminaorg/brzaguza/src/engines"

var Info engines.Info = engines.Info{
	Domain:         "www.bing.com",
	Name:           engines.BING,
	URL:            "https://www.bing.com/search?q=",
	ResultsPerPage: 10,
	Crawlers:       []engines.Name{engines.BING},
}

var dompaths engines.DOMPaths = engines.DOMPaths{
	Result:      "ol#b_results > li.b_algo",
	Link:        "h2 > a",
	Title:       "h2 > a",
	Description: "div.b_caption",
}

var Support engines.SupportedSettings = engines.SupportedSettings{}
