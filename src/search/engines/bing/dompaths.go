package bing

import (
	"github.com/hearchco/agent/src/search/scraper"
)

var dompaths = scraper.DOMPaths{
	Result:      "ol#b_results > li.b_algo",
	URL:         "h2 > a",
	Title:       "h2 > a",
	Description: "div.b_caption",
}
