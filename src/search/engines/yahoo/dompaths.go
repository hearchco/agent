package yahoo

import (
	"github.com/hearchco/agent/src/search/scraper"
)

var dompaths = scraper.DOMPaths{
	Result:      "div#main > div > div#web > ol > li > div.algo",
	URL:         "h3.title > a",
	Title:       "h3.title > a",
	Description: "div > div.compText > p > span",
}
