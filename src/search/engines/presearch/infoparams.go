package presearch

import (
	"github.com/hearchco/agent/src/search/engines"
	"github.com/hearchco/agent/src/search/scraper"
)

var info = scraper.Info{
	Name:    engines.PRESEARCH,
	Domain:  "presearch.com",
	URL:     "https://presearch.com/search",
	Origins: []engines.Name{engines.PRESEARCH, engines.GOOGLE},
}

var params = scraper.Params{
	Page:       "page",
	SafeSearch: "use_safe_search", // // Can be "true" or "false".
}
