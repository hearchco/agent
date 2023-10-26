package google

import (
	"github.com/tminaorg/brzaguza/src/engines"
)

var Info engines.Info = engines.Info{
	Domain:         "www.google.com",
	Name:           engines.GOOGLE,
	URL:            "https://www.google.com/search?q=",
	ResultsPerPage: 10,
	Crawlers:       []engines.Name{engines.GOOGLE},
}

/*
// This should be in Settings
var timings config.Timings = config.Timings{
	Timeout:     10 * time.Second, // the default in colly
	PageTimeout: 5 * time.Second,
	Delay:       100 * time.Millisecond,
	RandomDelay: 50 * time.Millisecond,
	Parallelism: 2, //two requests will be sent to the server, 100 + [0,50) milliseconds apart from the next two
}
*/

var dompaths engines.DOMPaths = engines.DOMPaths{
	Result:      "div.g",
	Link:        "a",
	Title:       "a > h3",
	Description: "div > span",
}

var Support engines.SupportedSettings = engines.SupportedSettings{}
