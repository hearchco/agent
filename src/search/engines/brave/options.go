package brave

import (
	"github.com/hearchco/hearchco/src/search/engines"
)

var Info engines.Info = engines.Info{
	Domain:         "search.brave.com",
	Name:           engines.BRAVE,
	URL:            "https://search.brave.com/search?q=",
	ResultsPerPage: 20,
}

var dompaths engines.DOMPaths = engines.DOMPaths{
	Result:      "div.snippet",
	Link:        "a",
	Title:       "div.title",
	Description: "div.snippet-description",
}

var Support engines.SupportedSettings = engines.SupportedSettings{
	Locale:     true,
	SafeSearch: true,
}
