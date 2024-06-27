package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

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

func writeResponseImageProxy(w http.ResponseWriter, resp *http.Response) error {
	w.Header().Set("Content-Encoding", resp.Header.Get("Content-Encoding"))
	w.Header().Set("Content-Length", resp.Header.Get("Content-Length"))
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.WriteHeader(resp.StatusCode)
	_, err := io.Copy(w, resp.Body)
	return err
}
