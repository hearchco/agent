package router

import (
	"net/http"
	"net/url"

	"github.com/rs/zerolog/log"
)

func writeResponse(w http.ResponseWriter, status int, body string) {
	w.WriteHeader(status)
	_, err := w.Write([]byte(body))
	if err != nil {
		log.Error().
			Err(err).
			Str("body", string(body)).
			Msg("Error writing response")
	}
}

func writeResponseJSON(w http.ResponseWriter, status int, body []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err := w.Write(body)
	if err != nil {
		log.Error().
			Err(err).
			Str("body", string(body)).
			Msg("Error writing response")
	}
}

func getParamOrDefault(params url.Values, key string, def ...string) string {
	val := params.Get(key)
	if val == "" && len(def) > 0 {
		return def[0]
	}
	return val
}
