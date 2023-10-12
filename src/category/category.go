package category

import (
	"strings"

	"github.com/rs/zerolog/log"
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

// returns category, rest of query
func FromQuery(query string) (Name, string) {
	if query[0] != '!' {
		return "", query
	}
	sp := strings.SplitN(query, " ", 2)
	cat := sp[0][1:]
	q := sp[1]
	val, ok := FromString[cat]
	if ok {
		return val, q
	}
	log.Debug().Msgf("invalid category in query: %v", query)
	return "", q
}
