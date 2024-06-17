package googlescholar

import (
	"github.com/hearchco/agent/src/search/scraper"
)

var dompaths = scraper.DOMPaths{
	Result:      "div#gs_res_ccl_mid > div.gs_or",
	URL:         "h3 > a",
	Title:       "h3 > a",
	Description: "div.gs_rs",
}
