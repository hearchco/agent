package presearch

type jsonResult struct {
	Title   string `json:"title"`
	Link    string `json:"link"`
	Desc    string `json:"description"`
	Favicon string `json:"favicon"`
}

type jsonResponse struct {
	Results struct {
		StandardResults []jsonResult `json:"standardResults"`
	} `json:"results"`
}
