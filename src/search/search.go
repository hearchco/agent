package search

import (
	"github.com/hearchco/hearchco/src/anonymize"
	"github.com/hearchco/hearchco/src/cache"
	"github.com/hearchco/hearchco/src/config"
	"github.com/hearchco/hearchco/src/search/engines"
	"github.com/hearchco/hearchco/src/search/result"
	"github.com/rs/zerolog/log"
)

func Search(query string, options engines.Options, conf config.Config, db cache.DB) ([]result.Result, bool) {
	var results []result.Result
	var foundInDB bool
	gerr := db.Get(query, &results)
	if gerr != nil {
		// Error in reading cache is not returned, just logged
		log.Error().
			Err(gerr).
			Str("queryAnon", anonymize.String(query)).
			Str("queryHash", anonymize.HashToSHA256B64(query)).
			Msg("cli.Run(): failed accessing cache")
	} else if results != nil {
		foundInDB = true
	} else {
		foundInDB = false
	}

	if foundInDB {
		log.Debug().
			Str("queryAnon", anonymize.String(query)).
			Str("queryHash", anonymize.HashToSHA256B64(query)).
			Msg("Found results in cache")
	} else {
		log.Debug().
			Str("queryAnon", anonymize.String(query)).
			Str("queryHash", anonymize.HashToSHA256B64(query)).
			Msg("Nothing found in cache, doing a clean search")
		results = PerformSearch(query, options, conf)
	}

	return results, foundInDB
}
