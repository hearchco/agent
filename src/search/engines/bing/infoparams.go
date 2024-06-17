package bing

import (
	"github.com/hearchco/agent/src/search/engines"
	"github.com/hearchco/agent/src/search/scraper"
)

var info = scraper.Info{
	Name:    engines.BING,
	Domain:  "www.bing.com",
	URL:     "https://www.bing.com/search",
	Origins: []engines.Name{engines.BING},
}

var params = scraper.Params{
	Page:       "first",
	Locale:     "setlang", // Should be first 2 characters of Locale.
	LocaleSec:  "cc",      // Should be last 2 characters of Locale.
	SafeSearch: "",        // Always enabled.
}
