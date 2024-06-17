package swisscows

import (
	"github.com/hearchco/agent/src/search/engines"
	"github.com/hearchco/agent/src/search/scraper"
)

var info = scraper.Info{
	Name:    engines.SWISSCOWS,
	Domain:  "swisscows.com",
	URL:     "https://api.swisscows.com/web/search",
	Origins: []engines.Name{engines.SWISSCOWS, engines.BING},
}

var params = scraper.Params{
	Page:   "offset",
	Locale: "region", // Should be the same as Locale, only with "_" replaced by "-".
}

const freshnessParam = "freshness=All"
const itemsParam = "itemsCount=10"
