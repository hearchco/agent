package yahoo

import (
	"github.com/hearchco/agent/src/search/engines"
	"github.com/hearchco/agent/src/search/scraper"
)

var info = scraper.Info{
	Name:    engines.YAHOO,
	Domain:  "search.yahoo.com",
	URL:     "https://search.yahoo.com/search",
	Origins: []engines.Name{engines.YAHOO, engines.BING},
}

var params = scraper.Params{
	Page:       "b",
	SafeSearch: "vm", // Can be "p" (disabled) or "r" (enabled).
}

const safeSearchCookiePrefix = "sB=v=1&pn=10&rw=new&userset=0"
