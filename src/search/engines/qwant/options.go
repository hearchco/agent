package qwant

import (
	"github.com/hearchco/hearchco/src/search/engines"
)

var Info = engines.Info{
	Domain:         "www.qwant.com",
	Name:           engines.QWANT,
	URL:            "https://api.qwant.com/v3/search/web?q=",
	ResultsPerPage: 10,
}

var Support = engines.SupportedSettings{
	Locale:                  true,
	SafeSearch:              true,
	Mobile:                  true,
	RequestedResultsPerPage: true,
}
