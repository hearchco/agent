package qwant

import "github.com/tminaorg/brzaguza/src/engines"

var Info engines.Info = engines.Info{
	Domain:         "www.qwant.com",
	Name:           "Qwant",
	URL:            "https://api.qwant.com/v3/search/web?q=",
	ResultsPerPage: 10,
	Crawlers:       []engines.Name{engines.Qwant, engines.Bing},
}

var Support engines.SupportedSettings = engines.SupportedSettings{
	Locale:                  true,
	SafeSearch:              true,
	Mobile:                  true,
	RequestedResultsPerPage: true,
}
