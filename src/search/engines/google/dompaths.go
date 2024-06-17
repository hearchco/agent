package google

import (
	"github.com/hearchco/agent/src/search/scraper"
)

var dompaths = scraper.DOMPaths{
	Result:      "div.g",
	URL:         "a",
	Title:       "a > h3",
	Description: "div > span",
}
