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

type bingImagesDomPaths struct {
	Result          string
	Metadata        string
	Title           string
	ImgFormatStr    string
	ThumbnailHeight string
	ThumbnailWidth  string
	Source          string
}

var dompaths = bingImagesDomPaths{
	Result:          "ul.dgControl_list > li",
	Metadata:        "a.iusc", // e.Attr("m")
	Title:           "div.infnmpt > div > ul > li > a",
	ImgFormatStr:    "div.imgpt > div > span",
	ThumbnailHeight: "a.iusc > div > img", // e.Attr("height")
	ThumbnailWidth:  "a.iusc > div > img", // e.Attr("width")
	Source:          "div.imgpt > div > div.lnkw > a",
}

var Support = engines.SupportedSettings{
	Locale: true,
}
