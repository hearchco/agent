package bingimages

import (
	"github.com/hearchco/agent/src/search/engines"
	"github.com/hearchco/agent/src/search/scraper"
)

var info = scraper.Info{
	Name:    engines.BINGIMAGES,
	Domain:  "www.bing.com",
	URL:     "https://www.bing.com/images/async",
	Origins: []engines.Name{engines.BINGIMAGES},
}

var params = scraper.Params{
	Page:       "first",
	Locale:     "setlang", // Should be first 2 characters of Locale.
	LocaleSec:  "cc",      // Should be last 2 characters of Locale.
	SafeSearch: "",        // Always enabled.
}

const asyncParam = "async=1"
const countParam = "count=35"
