package swisscows

type SCItem struct {
	Id         string `json:"id"`
	Title      string `json:"title"`
	Desc       string `json:"description"`
	URL        string `json:"url"`
	DisplayURL string `json:"displayUrl"`
}

type SCResponse struct {
	Items []SCItem `json:"items"`
}
