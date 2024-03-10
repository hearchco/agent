package googleimages

import (
	"github.com/hearchco/hearchco/src/search/engines"
)

var Info = engines.Info{
	Domain:         "images.google.com",
	Name:           engines.GOOGLEIMAGES,
	URL:            "https://www.google.com/search?q=",
	ResultsPerPage: 10,
}

var Support = engines.SupportedSettings{}
