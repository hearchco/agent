package etools

import (
	"github.com/hearchco/agent/src/search/scraper"
)

var dompaths = scraper.DOMPaths{
	Result:      "table.result > tbody > tr",
	URL:         "td.record > a",
	Title:       "td.record > a",
	Description: "td.record > div.text",
}
