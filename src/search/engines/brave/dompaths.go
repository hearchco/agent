package brave

import (
	"github.com/hearchco/agent/src/search/scraper"
)

var dompaths = scraper.DOMPaths{
	Result:      "div.snippet[data-type=\"web\"]",
	URL:         "a",
	Title:       "div.title",
	Description: "div.snippet-description",
}
