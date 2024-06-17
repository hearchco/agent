package result

type ImagesOutput struct {
	imagesOutputJSON
}

type imagesOutputJSON struct {
	Images

	URLHash          string `json:"url_hash,omitempty"`
	ThumbnailURLHash string `json:"thumbnail_url_hash,omitempty"`
}
