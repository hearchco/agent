package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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

func getParamOrDefault(params url.Values, key string, fallback ...string) string {
	val := params.Get(key)
	if val == "" && len(fallback) > 0 {
		return fallback[0]
	}
	return val
}
