package startpage

import (
	"github.com/hearchco/agent/src/search/scraper"
)

var dompaths = scraper.DOMPaths{
	Result:      "div.w-gl > div.result",
	URL:         "a.result-title",
	Title:       "a.result-title",
	Description: "p.description",
}
