package mojeek

import (
	"github.com/hearchco/agent/src/search/scraper"
)

var dompaths = scraper.DOMPaths{
	Result:      "ul.results-standard > li",
	URL:         "h2 > a.title",
	Title:       "h2 > a.title",
	Description: "p.s",
}
