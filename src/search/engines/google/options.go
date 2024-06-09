package google

import (
	"github.com/hearchco/hearchco/src/search/engines"
)

var Info = engines.Info{
	Domain:         "www.google.com",
	Name:           engines.GOOGLE,
	URL:            "https://google.com/search?hl=en&lr=lang_en&q=",
	ResultsPerPage: 10,
}

var dompaths = engines.DOMPaths{
	Result:      "div.g",
	Link:        "a",
	Title:       "a > h3",
	Description: "div > span",
}

var Support = engines.SupportedSettings{}
