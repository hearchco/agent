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
	if query[0] != '!' {
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
