package qwant

import "github.com/tminaorg/brzaguza/src/structures"

var Info structures.SEInfo = structures.SEInfo{
	Domain: "www.qwant.com",
	Name:   "Qwant",
	// not using "https://www.qwant.com/?q="
	URL:        "https://api.qwant.com/v3/search/web?q=",
	ResPerPage: 10,
	Crawlers:   []structures.EngineName{structures.Qwant, structures.Bing},
}
