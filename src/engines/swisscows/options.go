package swisscows

import "github.com/tminaorg/brzaguza/src/engines"

var Info engines.Info = engines.Info{
	Domain:         "swisscows.com",
	Name:           "Swisscows",
	URL:            "https://api.swisscows.com/web/search?",
	ResultsPerPage: 10,
	Crawlers:       []engines.Name{engines.Bing},
}

var Support engines.SupportedSettings = engines.SupportedSettings{
	Locale: true,
}
