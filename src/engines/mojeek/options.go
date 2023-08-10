package mojeek

import "github.com/tminaorg/brzaguza/src/structures"

var Info structures.SEInfo = structures.SEInfo{
	Domain:     "www.mojeek.com",
	Name:       "Mojeek",
	URL:        "https://www.mojeek.com/search?q=",
	ResPerPage: 10,
	Crawlers:   []structures.EngineName{structures.Mojeek},
}

var dompaths structures.SEDOMPaths = structures.SEDOMPaths{
	Result:      "ul.results-standard > li",
	Title:       "h2 > a.title",
	Description: "p.s",
}
