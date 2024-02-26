package search

import (
	"strings"

	"github.com/hearchco/hearchco/src/anonymize"
	"github.com/hearchco/hearchco/src/config"
	"github.com/hearchco/hearchco/src/search/category"
	"github.com/hearchco/hearchco/src/search/engines"
	"github.com/rs/zerolog/log"
)

func procBang(query string, options engines.Options, conf config.Config) (string, category.Name, config.Timings, []engines.Name) {
	useSpec, specEng := procSpecificEngine(query, conf)
	goodCat, cat := procCategory(query, options)
	if !goodCat && !useSpec && query[0] == '!' {
		// cat is set to GENERAL
		log.Debug().
			Str("queryAnon", anonymize.String(query)).
			Str("queryHash", anonymize.HashToSHA256B64(query)).
			Msg("search.procBang(): invalid bang (not category or engine shortcut)")
	}

	query = trimBang(query)

	if useSpec {
		return query, category.GENERAL, conf.Categories[category.GENERAL].Timings, []engines.Name{specEng}
	} else {
		return query, cat, conf.Categories[cat].Timings, conf.Categories[cat].Engines
	}
}

func trimBang(query string) string {
	if (query)[0] == '!' {
		return strings.SplitN(query, " ", 2)[1]
	}
	return query
}

func procSpecificEngine(query string, conf config.Config) (bool, engines.Name) {
	if query[0] != '!' {
		return false, engines.UNDEFINED
	}
	sp := strings.SplitN(query, " ", 2)
	bangWord := sp[0][1:]
	for key, val := range conf.Settings {
		if strings.EqualFold(bangWord, val.Shortcut) || strings.EqualFold(bangWord, key.String()) {
			return true, key
		}
	}

	return false, engines.UNDEFINED
}

func procCategory(query string, options engines.Options) (bool, category.Name) {
	cat := category.FromQuery(query)
	if cat != "" {
		return true, cat
	} else if options.Category == "" {
		return false, category.GENERAL
	} else {
		return false, options.Category
	}
}
