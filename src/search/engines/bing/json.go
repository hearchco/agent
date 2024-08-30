package bing

type imgJsonMetadata struct {
	PageURL      string `json:"purl"`
	ThumbnailURL string `json:"turl"`
	ImageURL     string `json:"murl"`
	Desc         string `json:"desc"`
}
