package brave

import "github.com/tminaorg/brzaguza/src/engines"

var Info engines.Info = engines.Info{
	Domain:         "search.brave.com",
	Name:           engines.Brave,
	URL:            "https://search.brave.com/search?q=",
	ResultsPerPage: 20,
}

var dompaths engines.DOMPaths = engines.DOMPaths{
	Result:      "div#results > div[class*=\"snippet fdb\"][data-type=\"web\"]",
	Link:        "a.result-header",
	Title:       "a.result-header > span.snippet-title",
	Description: "div.snippet-content > p.snippet-description",
}

var Support engines.SupportedSettings = engines.SupportedSettings{}
