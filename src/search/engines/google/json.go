package google

type imgJsonResponse struct {
	ISCHJ ischj `json:"ischj"`
}

type ischj struct {
	Metadata []metadata `json:"metadata"`
}

type metadata struct {
	Result        jsonResult `json:"result"`
	TextInGrid    textInGrid `json:"text_in_grid"`
	OriginalImage image      `json:"original_image"`
	Thumbnail     image      `json:"thumbnail"`
}

type jsonResult struct {
	ReferrerUrl string `json:"referrer_url"`
	PageTitle   string `json:"page_title"`
	SiteTitle   string `json:"site_title"`
}

type textInGrid struct {
	Snippet string `json:"snippet"`
}

type image struct {
	Url    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}
