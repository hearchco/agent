package google

import "github.com/tminaorg/brzaguza/src/structures"

var Info structures.SEInfo = structures.SEInfo{
	Domain:     "www.google.com",
	Name:       "Google",
	URL:        "https://www.google.com/search?q=",
	ResPerPage: 10,
	Crawlers:   []structures.EngineName{structures.Google},
}

var dompaths structures.SEDOMPaths = structures.SEDOMPaths{
	Result:      "div.g",
	Link:        "a",
	Title:       "div > div > div > a > h3",
	Description: "div > div > div > div:first-child > span:first-child",
}
