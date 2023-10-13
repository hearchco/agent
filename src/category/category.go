package category

import (
	"strings"
)

var FromString map[string]Name = map[string]Name{
	//main
	"general": GENERAL,
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
	val, ok := FromString[cat]
	if ok {
		return val
	}
	return ""
}
