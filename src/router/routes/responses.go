package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

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

func writeResponse(w http.ResponseWriter, status int, body string) error {
	w.WriteHeader(status)
	_, err := w.Write([]byte(body))
	return err
}

func writeResponseJSON(w http.ResponseWriter, status int, body any) error {
	res, err := json.Marshal(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, werr := w.Write([]byte("internal server error"))
		if werr != nil {
			return fmt.Errorf("%w: %w", werr, err)
		}
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(res)
	return err
}

func writeResponseSuggestions(w http.ResponseWriter, status int, query string, suggestions []string) error {
	jsonStruct := [...]any{query, suggestions}
	res, err := json.Marshal(jsonStruct)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, werr := w.Write([]byte("internal server error"))
		if werr != nil {
			return fmt.Errorf("%w: %w", werr, err)
		}
		return err
	}

	w.Header().Set("Content-Type", "application/x-suggestions+json")
	w.WriteHeader(status)
	_, err = w.Write(res)
	return err
}
