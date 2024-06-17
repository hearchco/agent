package startpage

import (
	"github.com/hearchco/agent/src/search/engines"
	"github.com/hearchco/agent/src/search/scraper"
)

var info = scraper.Info{
	Name:    engines.STARTPAGE,
	Domain:  "www.startpage.com",
	URL:     "https://www.startpage.com/sp/search",
	Origins: []engines.Name{engines.STARTPAGE, engines.GOOGLE},
}

var params = scraper.Params{
	Page:       "page",
	SafeSearch: "qadf", // Can be "none" or empty param (empty means it's enabled).
}
