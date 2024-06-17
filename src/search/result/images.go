package result

import (
	"github.com/hearchco/agent/src/utils/anonymize"
)

type Images struct {
	imagesJSON
}

type imagesJSON struct {
	General

	OriginalSize  ImageFormat `json:"original"`
	ThumbnailSize ImageFormat `json:"thumbnail"`
	ThumbnailURL  string      `json:"thumbnail_url"`
	SourceName    string      `json:"source"`
	SourceURL     string      `json:"source_url"`
}

type ImageFormat struct {
	Height int `json:"height"`
	Width  int `json:"width"`
}

func (r Images) OriginalSize() ImageFormat {
	return r.imagesJSON.OriginalSize
}

func (r Images) ThumbnailSize() ImageFormat {
	return r.imagesJSON.ThumbnailSize
}

func (r Images) ThumbnailURL() string {
	return r.imagesJSON.ThumbnailURL
}

func (r Images) SourceName() string {
	return r.imagesJSON.SourceName
}

func (r Images) SourceURL() string {
	return r.imagesJSON.SourceURL
}

func (r Images) ConvertToOutput(salt string) ResultOutput {
	return ImagesOutput{
		imagesOutputJSON{
			r,
			anonymize.HashToSHA256B64Salted(r.URL(), salt),
			anonymize.HashToSHA256B64Salted(r.ThumbnailURL(), salt),
		},
	}
}
