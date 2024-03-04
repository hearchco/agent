package bingimages

import (
	"github.com/hearchco/hearchco/src/search/engines"
)

var params = []string{"&async=1", "&count=35"}

var Info = engines.Info{
	Domain:         "www.bing.com",
	Name:           engines.BINGIMAGES,
	URL:            "https://www.bing.com/images/async?q=",
	ResultsPerPage: 35,
}

var Support = engines.SupportedSettings{
	Locale: true,
}
