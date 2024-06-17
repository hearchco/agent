package qwant

import (
	"github.com/hearchco/agent/src/search/engines"
	"github.com/hearchco/agent/src/search/scraper"
)

var info = scraper.Info{
	Name:    engines.QWANT,
	Domain:  "www.qwant.com",
	URL:     "https://api.qwant.com/v3/search/web",
	Origins: []engines.Name{engines.QWANT, engines.BING},
}

var params = scraper.Params{
	Page:       "offset",
	Locale:     "locale",     // Same as Locale, only the last two characters are lowered and not everything is supported.
	SafeSearch: "safesearch", // Can be "0" or "1".
}

const countParam = "count=10"
