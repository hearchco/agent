package brave

import (
	"github.com/hearchco/agent/src/search/engines"
	"github.com/hearchco/agent/src/search/scraper"
)

var info = scraper.Info{
	Name:    engines.BRAVE,
	Domain:  "search.brave.com",
	URL:     "https://search.brave.com/search",
	Origins: []engines.Name{engines.BRAVE, engines.GOOGLE},
}

var params = scraper.Params{
	Page:       "offset",
	Locale:     "country",    // Should be last 2 characters of Locale.
	SafeSearch: "safesearch", // Can be "off" or "strict".
}

const sourceParam = "source=web"
const spellcheckParam = "spellcheck=0"
