package google

import (
	"github.com/hearchco/hearchco/src/search/engines"
)

var Info engines.Info = engines.Info{
	Domain:         "www.google.com",
	Name:           engines.GOOGLE,
	URL:            "https://www.google.com/search?q=",
	ResultsPerPage: 10,
}

var dompaths engines.DOMPaths = engines.DOMPaths{
	Result:      "div.g",
	Link:        "a",
	Title:       "a > h3",
	Description: "div > span",
}

var Support engines.SupportedSettings = engines.SupportedSettings{}
