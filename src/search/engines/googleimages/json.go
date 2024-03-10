package googleimages

type Result struct {
	ReferrerUrl string `json:"referrer_url"`
	PageTitle   string `json:"page_title"`
	SiteTitle   string `json:"site_title"`
}

type TextInGrid struct {
	Snippet string `json:"snippet"`
}

type Image struct {
	Url    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

type Metadata struct {
	Result        Result     `json:"result"`
	TextInGrid    TextInGrid `json:"text_in_grid"`
	OriginalImage Image      `json:"original_image"`
	Thumbnail     Image      `json:"thumbnail"`
}

type ISCHJ struct {
	Metadata []Metadata `json:"metadata"`
}

type JsonResponse struct {
	ISCHJ ISCHJ `json:"ischj"`
}
