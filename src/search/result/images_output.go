package result

type ImagesOutput struct {
	imagesOutputJSON
}

type imagesOutputJSON struct {
	Images

	URLHash                   string `json:"url_hash,omitempty"`
	URLHashTimestamp          string `json:"url_hash_timestamp,omitempty"`
	ThumbnailURLHash          string `json:"thumbnail_url_hash,omitempty"`
	ThumbnailURLHashTimestamp string `json:"thumbnail_url_hash_timestamp,omitempty"`
}
