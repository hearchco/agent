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
		log.Error().Err(err1).Msgf("%v: Failed body unmarshall to json:\n%v\n", Info.Name, string(body))
		return nil
	}
	byteMain, err2 := json.Marshal(yr[1])
	if err2 != nil {
		log.Error().Err(err2).Msgf("%v: Failed marshalling the relevant json content:\n%v\n", Info.Name, string(body))
		return nil
	}

	var mainContent MainContent
	err3 := json.Unmarshal(byteMain, &mainContent)
	if err3 != nil {
		log.Error().Err(err3).Msgf("%v: Failed unmarshalling to MainContent:\n%v\n", Info.Name, string(byteMain))
		return nil
	}

	return &mainContent
}
