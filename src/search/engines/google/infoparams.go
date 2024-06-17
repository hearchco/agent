package google

import (
	"github.com/hearchco/agent/src/search/engines"
	"github.com/hearchco/agent/src/search/scraper"
)

var info = scraper.Info{
	Name:    engines.GOOGLE,
	Domain:  "www.google.com",
	URL:     "https://www.google.com/search",
	Origins: []engines.Name{engines.GOOGLE},
}

var params = scraper.Params{
	Page:       "start",
	Locale:     "hl",   // Should be first 2 characters of Locale.
	LocaleSec:  "lr",   // Should be first 2 characters of Locale with prefixed "lang_".
	SafeSearch: "safe", // Can be "off", "medium or "high".
}

const filterParam = "filter=0"
