package mojeek

import (
	"github.com/hearchco/agent/src/search/engines"
	"github.com/hearchco/agent/src/search/scraper"
)

var info = scraper.Info{
	Name:    engines.MOJEEK,
	Domain:  "www.mojeek.com",
	URL:     "https://www.mojeek.com/search",
	Origins: []engines.Name{engines.MOJEEK},
}

var params = scraper.Params{
	Page:       "s",
	Locale:     "lb",   // Should be first 2 characters of Locale.
	LocaleSec:  "arc",  // Should be last 2 characters of Locale.
	SafeSearch: "safe", // Can be "0" or "1".
}
