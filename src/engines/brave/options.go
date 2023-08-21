package brave

import "github.com/tminaorg/brzaguza/src/structures"

var Info structures.SEInfo = structures.SEInfo{
	Domain:         "search.brave.com",
	Name:           "Brave",
	URL:            "https://search.brave.com/search?q=",
	ResultsPerPage: 20,
}

var dompaths structures.SEDOMPaths = structures.SEDOMPaths{
	Result:      "div#results > div[class*=\"snippet fdb\"][data-type=\"web\"]",
	Link:        "a.result-header",
	Title:       "a.result-header > span.snippet-title",
	Description: "div.snippet-content > p.snippet-description",
}

var Support structures.SupportedSettings = structures.SupportedSettings{}
