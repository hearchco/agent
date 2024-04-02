package category

import (
	"strings"
)

var FromString = map[string]Name{
	//main
	"general": GENERAL,
	"images":  IMAGES,
	"info":    INFO,
	"science": SCIENCE,
	"news":    NEWS,
	"blog":    BLOG,
	"surf":    SURF,
	"newnews": NEWNEWS,
	//alternatives
	"wiki":  INFO,
	"sci":   SCIENCE,
	"nnews": NEWNEWS,
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

// Returns category from query or fallback or GENERAL and whether the category is from the query or not
func FromQueryWithFallback(query string, fallback Name) (Name, bool) {
	cat := FromQuery(query)
	if cat != "" {
		return cat, true
	} else if fallback != "" {
		return fallback, false
	} else {
		return GENERAL, false
	}
}
