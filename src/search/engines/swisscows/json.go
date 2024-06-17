package swisscows

type jsonResponse struct {
	Items []jsonItem `json:"items"`
}

type jsonItem struct {
	Id         string `json:"id"`
	Title      string `json:"title"`
	Desc       string `json:"description"`
	URL        string `json:"url"`
	DisplayURL string `json:"displayUrl"`
}
