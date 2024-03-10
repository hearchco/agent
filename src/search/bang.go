package search

import (
	"strings"

	"github.com/hearchco/hearchco/src/anonymize"
	"github.com/hearchco/hearchco/src/config"
	"github.com/hearchco/hearchco/src/search/category"
	"github.com/hearchco/hearchco/src/search/engines"
	"github.com/rs/zerolog/log"
)

func procBang(query string, setCategory category.Name, settings map[engines.Name]config.Settings, categories map[category.Name]config.Category) (string, category.Name, config.Timings, []engines.Name) {
	useSpec, specEng := procSpecificEngine(query, settings)
	goodCat, cat := procCategory(query, setCategory)
	if !goodCat && !useSpec && (query != "" && query[0] == '!') {
		// cat is set to GENERAL
		log.Debug().
			Str("queryAnon", anonymize.String(query)).
			Str("queryHash", anonymize.HashToSHA256B64(query)).
			Msg("search.procBang(): invalid bang (not category or engine shortcut)")
	}

	query = trimBang(query)

	if useSpec {
		return query, category.GENERAL, categories[category.GENERAL].Timings, []engines.Name{specEng}
	} else {
		return query, cat, categories[cat].Timings, categories[cat].Engines
	}
}

// takes the bang out of the query performs TrimSpace
func trimBang(query string) string {
	query = strings.TrimSpace(query)

	if query == "" || query[0] != '!' {
		return query
	}

	sp := strings.SplitN(query, " ", 2)
	if len(sp) == 1 {
		// only the bang is present
		return ""
	}

	return strings.TrimSpace(sp[1])
}

func procSpecificEngine(query string, settings map[engines.Name]config.Settings) (bool, engines.Name) {
	if query == "" || query[0] != '!' {
		return false, engines.UNDEFINED
	}
	sp := strings.SplitN(query, " ", 2)
	bangWord := sp[0][1:]
	for key, val := range settings {
		if strings.EqualFold(bangWord, val.Shortcut) || strings.EqualFold(bangWord, key.String()) {
			return true, key
		}
	}

	return false, engines.UNDEFINED
}

// returns category in the query if a valid category is present
func procCategory(query string, setCategory category.Name) (bool, category.Name) {
	cat := category.FromQuery(query)
	if cat != "" {
		return true, cat
	} else if setCategory == "" {
		return false, category.GENERAL
	} else {
		return false, setCategory
	}
}
