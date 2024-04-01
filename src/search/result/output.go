package result

type GeneralResultOutput struct {
	URL         string          `json:"url"`
	URLHash     string          `json:"url_hash,omitempty"`
	Rank        uint            `json:"rank"`
	Score       float64         `json:"score"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	EngineRanks []RetrievedRank `json:"engine_ranks"`
}

func ConvertToGeneralOutput(results []Result) []GeneralResultOutput {
	resultsOutput := make([]GeneralResultOutput, 0, len(results))
	for _, r := range results {
		resultsOutput = append(resultsOutput, GeneralResultOutput{
			URL:         r.URL,
			URLHash:     r.URLHash,
			Rank:        r.Rank,
			Score:       r.Score,
			Title:       r.Title,
			Description: r.Description,
			EngineRanks: r.EngineRanks,
		})
	}
	return resultsOutput
}

type ImageResultOutput struct {
	URL         string          `json:"url"`
	URLHash     string          `json:"url_hash,omitempty"`
	Rank        uint            `json:"rank"`
	Score       float64         `json:"score"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	EngineRanks []RetrievedRank `json:"engine_ranks"`
	ImageResult ImageResult     `json:"image_result"`
}

func ConvertToImageOutput(results []Result) []ImageResultOutput {
	resultsOutput := make([]ImageResultOutput, 0, len(results))
	for _, r := range results {
		resultsOutput = append(resultsOutput, ImageResultOutput{
			URL:         r.URL,
			URLHash:     r.URLHash,
			Rank:        r.Rank,
			Score:       r.Score,
			Title:       r.Title,
			Description: r.Description,
			EngineRanks: r.EngineRanks,
			ImageResult: r.ImageResult,
		})
	}
	return resultsOutput
}
