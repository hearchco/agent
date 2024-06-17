package duckduckgo

import (
	"github.com/hearchco/agent/src/search/engines"
	"github.com/hearchco/agent/src/search/scraper"
)

var info = scraper.Info{
	Name:    engines.DUCKDUCKGO,
	Domain:  "lite.duckduckgo.com",
	URL:     "https://lite.duckduckgo.com/lite/",
	Origins: []engines.Name{engines.DUCKDUCKGO, engines.BING},
}

var params = scraper.Params{
	Page:       "dc",
	Locale:     "kl", // Should be Locale with _ replaced by - and first 2 letters as last and vice versa.
	SafeSearch: "",   // Always enabled.
}
