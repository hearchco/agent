package googleimages

import (
	"github.com/hearchco/agent/src/search/engines"
	"github.com/hearchco/agent/src/search/scraper"
)

var info = scraper.Info{
	Name:    engines.GOOGLEIMAGES,
	Domain:  "images.google.com",
	URL:     "https://www.google.com/search",
	Origins: []engines.Name{engines.GOOGLEIMAGES},
}

var params = scraper.Params{
	Page:       "async=_fmt:json,p:1,ijn",
	Locale:     "hl",   // Should be first 2 characters of Locale.
	LocaleSec:  "lr",   // Should be first 2 characters of Locale with prefixed "lang_".
	SafeSearch: "safe", // Can be "off", "medium or "high".
}

const tbmParam = "tbm=isch"
const asearchParam = "asearch=isch"
const filterParam = "filter=0"
