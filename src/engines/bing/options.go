package bing

import "github.com/tminaorg/brzaguza/src/structures"

var Info structures.SEInfo = structures.SEInfo{
	Domain:         "www.bing.com",
	Name:           "Bing",
	URL:            "https://www.bing.com/search?q=",
	ResultsPerPage: 10,
	Crawlers:       []structures.EngineName{structures.Brave, structures.Google},
}

var dompaths structures.SEDOMPaths = structures.SEDOMPaths{
	Result:      "ol#b_results > li.b_algo",
	Link:        "h2 > a",
	Title:       "h2 > a",
	Description: "div.b_caption",
}

var Support structures.SupportedSettings = structures.SupportedSettings{}
