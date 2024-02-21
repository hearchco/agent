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

type thumbnailDomPaths struct {
	Path   string
	Height string
	Width  string
}

type bingImagesDomPaths struct {
	Result       string
	Metadata     string
	Title        string
	ImgFormatStr string
	Thumbnail    []thumbnailDomPaths
	Source       string
}

var dompaths = bingImagesDomPaths{
	Result:       "ul.dgControl_list > li",
	Metadata:     "a.iusc", // e.Attr("m")
	Title:        "div.infnmpt > div > ul > li > a",
	ImgFormatStr: "div.imgpt > div > span",
	Thumbnail: []thumbnailDomPaths{
		{
			Path:   "a.iusc > div > img.mimg",
			Height: "height",
			Width:  "width",
		},
		{
			Path:   "a.iusc > div > div > div.mimg > div",
			Height: "data-height",
			Width:  "data-width",
		},
		{
			Path:   "a.iusc > div > div > div.mimg > img",
			Height: "height",
			Width:  "width",
		},
	},
	Source: "div.imgpt > div > div.lnkw > a",
}

var Support = engines.SupportedSettings{
	Locale: true,
}
