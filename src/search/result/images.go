package result

import (
	"github.com/rs/zerolog/log"

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
	if r.imagesJSON.OriginalSize.Height == 0 || r.imagesJSON.OriginalSize.Width == 0 {
		log.Panic().
			Int("height", r.imagesJSON.OriginalSize.Height).
			Int("width", r.imagesJSON.OriginalSize.Width).
			Msg("OriginalSize is zero")
		// ^PANIC - Assert because the OriginalSize should never be zero.
	}

	return r.imagesJSON.OriginalSize
}

func (r Images) ThumbnailSize() ImageFormat {
	if r.imagesJSON.ThumbnailSize.Height == 0 || r.imagesJSON.ThumbnailSize.Width == 0 {
		log.Panic().
			Int("height", r.imagesJSON.ThumbnailSize.Height).
			Int("width", r.imagesJSON.ThumbnailSize.Width).
			Msg("ThumbnailSize is zero")
		// ^PANIC - Assert because the ThumbnailSize should never be zero.
	}

	return r.imagesJSON.ThumbnailSize
}

func (r Images) ThumbnailURL() string {
	if r.imagesJSON.ThumbnailURL == "" {
		log.Panic().Msg("ThumbnailURL is empty")
		// ^PANIC - Assert because the ThumbnailURL should never be empty.
	}

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
			anonymize.CalculateMACBase64(r.URL(), salt),
			anonymize.CalculateMACBase64(r.ThumbnailURL(), salt),
		},
	}
}
