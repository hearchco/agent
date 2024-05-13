package category

import (
	"strings"
)

var FromString = map[string]Name{
	"general":  GENERAL,
	"images":   IMAGES,
	"science":  SCIENCE,
	"sci":      SCIENCE,
	"quick":    QUICK,
	"fast":     QUICK,
	"thorough": THOROUGH,
	"slow":     THOROUGH,
}

// returns category
func FromQuery(query string) Name {
	if query == "" || query[0] != '!' {
		return ""
	}
	cat := strings.SplitN(query, " ", 2)[0][1:]
	if val, ok := FromString[cat]; ok {
		return val
	}
	return ""
}

func SafeFromString(cat string) Name {
	if cat == "" {
		return ""
	}
	ret, ok := FromString[cat]
	if !ok {
		return UNDEFINED
	}
	return ret
}

func FromQueryWithFallback(query string, fallback Name) Name {
	cat := FromQuery(query)
	if cat != "" {
		return cat
	} else if fallback != "" {
		return fallback
	} else {
		return GENERAL
	}
}
