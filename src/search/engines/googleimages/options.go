package googleimages

import "github.com/hearchco/hearchco/src/search/engines"

var params = "&tbm=isch&asearch=isch&async=_fmt:json,p:1,ijn:"

var Info engines.Info = engines.Info{
	Domain:         "images.google.com",
	Name:           engines.GOOGLEIMAGES,
	URL:            "https://www.google.com/search?q=",
	ResultsPerPage: 50,
}

var Support engines.SupportedSettings = engines.SupportedSettings{}
