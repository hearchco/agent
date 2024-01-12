package yep

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
)

type YepResponse []interface{}

type Result struct {
	URL       string `json:"url"`
	Title     string `json:"title"`
	TType     string `json:"type"`
	Snippet   string `json:"snippet"`
	VisualURL string `json:"visual_url"`
	FirstSeen string `json:"first_seen"`
}

type MainContent struct {
	Total   int      `json:"total"`
	Results []Result `json:"results"`
}

// This is pretty inefficient, and can probably be done better
func parseJSON(body []byte) *MainContent {
	var yr YepResponse

	err1 := json.Unmarshal(body, &yr)
	if err1 != nil {
		log.Error().
			Err(err1).
			Str("SEName", Info.Name.String()).
			Str("Body", string(body)).
			Msg("Failed body unmarshall to json")
		return nil
	}

	byteMain, err2 := json.Marshal(yr[1])
	if err2 != nil {
		log.Error().
			Err(err2).
			Str("SEName", Info.Name.String()).
			Str("Body", string(body)).
			Msg("Failed marshalling the relevant json content")
		return nil
	}

	var mainContent MainContent
	err3 := json.Unmarshal(byteMain, &mainContent)
	if err3 != nil {
		log.Error().
			Err(err3).
			Str("SEName", Info.Name.String()).
			Str("Body", string(body)).
			Msg("Failed unmarshalling to MainContent")
		return nil
	}

	return &mainContent
}
