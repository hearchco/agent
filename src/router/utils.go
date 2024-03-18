package router

import (
	"net/http"
	"net/url"
)

func writeResponse(w http.ResponseWriter, status int, body string) {
	w.WriteHeader(status)
	w.Write([]byte(body))
}

func writeResponseJSON(w http.ResponseWriter, status int, body []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(body)
}

func getParamOrDefault(params url.Values, key string, def ...string) string {
	val := params.Get(key)
	if val == "" && len(def) > 0 {
		return def[0]
	}
	return val
}
