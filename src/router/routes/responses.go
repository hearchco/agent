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

type ResultsResponse struct {
	Version  string                `json:"version"`
	Duration int64                 `json:"duration"`
	Results  []result.ResultOutput `json:"results"`
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

func writeResponseSuggestions(w http.ResponseWriter, status int, query string, suggestions []string, api bool) error {
	jsonStruct := make([]any, 2)
	jsonStruct[0] = query
	jsonStruct[1] = suggestions

	res, err := json.Marshal(jsonStruct)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, werr := w.Write([]byte("internal server error"))
		if werr != nil {
			return fmt.Errorf("%w: %w", werr, err)
		}
		return err
	}

	if api {
		w.Header().Set("Content-Type", "application/x-suggestions+json")
	} else {
		w.Header().Set("Content-Type", "application/json")
	}
	w.WriteHeader(status)
	_, err = w.Write(res)
	return err
}
