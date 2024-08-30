package bing

import (
	"github.com/hearchco/agent/src/search/scraper"
)

var dompaths = scraper.DOMPaths{
	Result:      "ol#b_results > li.b_algo",
	URL:         "h2 > a",
	Title:       "h2 > a",
	Description: "div.b_caption",
}

type thumbnailDomPaths struct {
	Path   string
	Height string
	Width  string
}

type metadataDomPaths struct {
	Path string
	Attr string
}

type bingImagesDomPaths struct {
	Result       string
	Metadata     metadataDomPaths
	Title        string
	ImgFormatStr string
	Thumbnail    [3]thumbnailDomPaths
	Source       string
}

var imageDompaths = bingImagesDomPaths{
	// aria-live is also a possible attribute for not()
	Result: "ul.dgControl_list > li[data-idx] > div.iuscp:not([vrhatt])",
	Metadata: metadataDomPaths{
		Path: "a.iusc",
		Attr: "m",
	},
	Title:        "div.infnmpt > div > ul > li > a",
	ImgFormatStr: "div.imgpt > div > span",
	Thumbnail: [...]thumbnailDomPaths{
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
	Source: "div.imgpt > div.img_info > div.lnkw > a",
}
