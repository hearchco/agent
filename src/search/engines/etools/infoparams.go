package etools

import (
	"github.com/hearchco/agent/src/search/engines"
	"github.com/hearchco/agent/src/search/scraper"
)

var info = scraper.Info{
	Name:    engines.ETOOLS,
	Domain:  "www.etools.ch",
	URL:     "https://www.etools.ch/searchSubmit.do",
	Origins: []engines.Name{engines.ETOOLS}, // Disabled because ETOOLS has issues most of the time: []engines.Name{engines.BING, engines.BRAVE, engines.DUCKDUCKGO, engines.GOOGLE, engines.MOJEEK, engines.QWANT, engines.YAHOO},
}

const pageURL = "https://www.etools.ch/search.do"

var params = scraper.Params{
	Page:       "page",
	SafeSearch: "safeSearch", // Can be "true" or "false".
}

const countryParam = "country=web"
const languageParam = "language=all"
