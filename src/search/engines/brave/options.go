package brave

import (
	"github.com/hearchco/hearchco/src/search/engines"
)

var Info = engines.Info{
	Domain:         "search.brave.com",
	Name:           engines.BRAVE,
	URL:            "https://search.brave.com/search?q=",
	ResultsPerPage: 20,
}

var dompaths = engines.DOMPaths{
	Result:      "div.snippet[data-type=\"web\"]",
	Link:        "a",
	Title:       "div.title",
	Description: "div.snippet-description",
}

var Support = engines.SupportedSettings{
	Locale:     true,
	SafeSearch: true,
}
