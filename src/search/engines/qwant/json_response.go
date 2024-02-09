package qwant

type QwantResults struct {
	Title       string `json:"title"`
	URL         string `json:"url"` //there is also a source field, what is it?
	Description string `json:"desc"`
}

type QwantMainlineItems struct {
	Type  string         `json:"type"`
	Items []QwantResults `json:"items"`
}

type QwantResponse struct {
	Status string `json:"status"`
	Data   struct {
		Res struct {
			Items struct {
				Mainline []QwantMainlineItems `json:"mainline"`
			} `json:"items"`
		} `json:"result"`
	} `json:"data"`
}
