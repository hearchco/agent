package presearch

type Result struct {
	Title   string `json:"title"`
	Link    string `json:"link"`
	Desc    string `json:"description"`
	Favicon string `json:"favicon"`
}

type PresearchResponse struct {
	Results struct {
		StandardResults []Result `json:"standardResults"`
	} `json:"results"`
}
