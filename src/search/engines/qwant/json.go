package qwant

type jsonResponse struct {
	Status string `json:"status"`
	Data   struct {
		Res struct {
			Items struct {
				Mainline []jsonMainlineItems `json:"mainline"`
			} `json:"items"`
		} `json:"result"`
	} `json:"data"`
}

type jsonMainlineItems struct {
	Type  string        `json:"type"`
	Items []jsonResults `json:"items"`
}

type jsonResults struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	Description string `json:"desc"`
}
