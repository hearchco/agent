package yahoo

import "github.com/hearchco/hearchco/src/engines"

// doesn't catch the yt videos
// the title cathes the link - e.g.: teentitans.fandom.com › wiki › Nya-NyaNya-Nya | Teen Titans Wiki | Fandom
// but should be just: Nya-Nya | Teen Titans Wiki | Fandom

var Info engines.Info = engines.Info{
	Domain:         "search.yahoo.com",
	Name:           engines.YAHOO,
	URL:            "https://search.yahoo.com/search?p=",
	ResultsPerPage: 10,
	Crawlers:       []engines.Name{engines.BING},
}

var dompaths engines.DOMPaths = engines.DOMPaths{
	Result:      "div#main > div > div#web > ol > li",
	Link:        "h3.title > a",
	Title:       "h3.title > a",
	Description: "div > div.compText > p > span",
}

var Support engines.SupportedSettings = engines.SupportedSettings{}
