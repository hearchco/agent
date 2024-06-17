package duckduckgo

import (
	"github.com/hearchco/agent/src/search/scraper"
)

var dompaths = scraper.DOMPaths{
	ResultsContainer: "div.filters > table > tbody",
	URL:              "td > a.result-link",
	Title:            "td > a.result-link",
	Description:      "td.result-snippet",
}
