package swisscows

import "github.com/tminaorg/brzaguza/src/structures"

var Info structures.SEInfo = structures.SEInfo{
	Domain:     "swisscows.com",
	Name:       "Swisscows",
	URL:        "https://api.swisscows.com/web/search?",
	ResPerPage: 10,
	Crawlers:   []structures.EngineName{structures.Bing},
}
