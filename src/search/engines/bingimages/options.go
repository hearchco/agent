package bingimages

import "github.com/hearchco/hearchco/src/search/engines"

var params = []string{"&async=1", "&count=35"}

var Info = engines.Info{
	Domain:         "www.bing.com",
	Name:           engines.BINGIMAGES,
	URL:            "https://www.bing.com/images/async?q=",
	ResultsPerPage: 10,
	Crawlers:       []engines.Name{engines.BINGIMAGES},
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
	Title:           "div.infnmpt > a",
	ImgFormatStr:    "div.imgpt > div > span",
	ThumbnailHeight: "a.iusc > img", // e.Attr("height")
	ThumbnailWidth:  "a.iusc > img", // e.Attr("width")
	Source:          "div.imgpt > div.lnkw > a",
}

var Support = engines.SupportedSettings{
	Locale: true,
}
