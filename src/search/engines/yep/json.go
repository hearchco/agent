package yep

type Result struct {
	URL     string `json:"url"`
	Title   string `json:"title"`
	TType   string `json:"type"`
	Snippet string `json:"snippet"`
	// VisualURL string `json:"visual_url"`
	// FirstSeen string `json:"first_seen"`
}

type JsonResponse struct {
	// Total   int      `json:"total"`
	Results []Result `json:"results"`
}
