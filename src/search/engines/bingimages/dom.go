package bingimages

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
	Thumbnail    []thumbnailDomPaths
	Source       string
}

var dompaths = bingImagesDomPaths{
	Result: "ul.dgControl_list > li[data-idx]",
	Metadata: metadataDomPaths{
		Path: "a.iusc",
		Attr: "m",
	},
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
