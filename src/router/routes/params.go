package routes

import (
	"net/url"
)

func getParamOrDefault(params url.Values, key string, fallback ...string) string {
	val := params.Get(key)
	if val == "" && len(fallback) > 0 {
		return fallback[0]
	}
	return val
}
