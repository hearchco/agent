package qwant

import "github.com/tminaorg/brzaguza/src/structures"

var Info structures.SEInfo = structures.SEInfo{
	Domain:         "www.qwant.com",
	Name:           "Qwant",
	URL:            "https://api.qwant.com/v3/search/web?q=",
	ResultsPerPage: 10,
	Crawlers:       []structures.EngineName{structures.Qwant, structures.Bing},
}

var Support structures.SupportedSettings = structures.SupportedSettings{
	Locale:                  true,
	SafeSearch:              true,
	Mobile:                  true,
	RequestedResultsPerPage: true,
}
