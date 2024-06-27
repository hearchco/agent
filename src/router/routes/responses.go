package routes

import (
	"github.com/hearchco/agent/src/search/result"
)

type ErrorResponse struct {
	Message string `json:"message"`
	Value   string `json:"value"`
}

type responseBase struct {
	Version  string `json:"version"`
	Duration int64  `json:"duration"`
}

type ResultsResponse struct {
	responseBase

	Results []result.ResultOutput `json:"results"`
}

type SuggestionsResponse struct {
	responseBase

	Suggestions []result.Suggestion `json:"suggestions"`
}
