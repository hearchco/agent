package search

import (
	"github.com/hearchco/hearchco/src/anonymize"
	"github.com/hearchco/hearchco/src/cache"
	"github.com/hearchco/hearchco/src/config"
	"github.com/hearchco/hearchco/src/search/engines"
	"github.com/hearchco/hearchco/src/search/result"
	"github.com/rs/zerolog/log"
)

func Search(query string, options engines.Options, db cache.DB, categoryConf config.Category, settings map[engines.Name]config.Settings, salt string) ([]result.Result, bool) {
	if results, err := db.GetResults(query, options.Category); err != nil {
		// Error in reading cache is not returned, just logged
		log.Error().
			Caller().
			Err(err).
			Str("queryAnon", anonymize.String(query)).
			Str("queryHash", anonymize.HashToSHA256B64(query)).
			Msg("Failed accessing cache")
	} else if results != nil {
		log.Debug().
			Str("queryAnon", anonymize.String(query)).
			Str("queryHash", anonymize.HashToSHA256B64(query)).
			Msg("Found results in cache")

		return results, true
	}

	// if the cache is inaccesible or the query+category is not in the cache
	log.Debug().
		Str("queryAnon", anonymize.String(query)).
		Str("queryHash", anonymize.HashToSHA256B64(query)).
		Msg("Nothing found in cache, doing a clean search")

	// the main line
	results := PerformSearch(query, options, categoryConf, settings, salt)
	result.Shorten(results, 2500)

	return results, false
}
