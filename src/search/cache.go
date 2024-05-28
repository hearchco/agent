package search

import (
	"github.com/hearchco/hearchco/src/anonymize"
	"github.com/hearchco/hearchco/src/cache"
	"github.com/hearchco/hearchco/src/config"
	"github.com/hearchco/hearchco/src/search/engines"
	"github.com/hearchco/hearchco/src/search/result"
	"github.com/rs/zerolog/log"
)

func CacheAndUpdateResults(
	query string, options engines.Options, db cache.DB,
	ttlConf config.TTL, categoryConf config.Category, settings map[engines.Name]config.Settings,
	results []result.Result, foundInDB bool,
	salt string,
) {
	if !foundInDB {
		log.Debug().
			Str("queryAnon", anonymize.String(query)).
			Str("queryHash", anonymize.HashToSHA256B64(query)).
			Msg("Caching results...")
		serr := db.SetResults(query, options.Category, results, ttlConf.Time)
		if serr != nil {
			log.Error().
				Caller().
				Err(serr).
				Str("queryAnon", anonymize.String(query)).
				Str("queryHash", anonymize.HashToSHA256B64(query)).
				Msg("Error updating database with search results")
		}
	} else {
		log.Debug().
			Str("queryAnon", anonymize.String(query)).
			Str("queryHash", anonymize.HashToSHA256B64(query)).
			Msg("Checking if results need to be updated")
		ttl, terr := db.GetResultsTTL(query, options.Category)
		if terr != nil {
			log.Error().
				Caller().
				Err(terr).
				Str("queryAnon", anonymize.String(query)).
				Str("queryHash", anonymize.HashToSHA256B64(query)).
				Msg("Error getting TTL from database")
		} else if ttl < ttlConf.RefreshTime {
			log.Info().
				Str("queryAnon", anonymize.String(query)).
				Str("queryHash", anonymize.HashToSHA256B64(query)).
				Msg("Updating results...")
			newResults := PerformSearch(query, options, categoryConf, settings, salt)
			uerr := db.SetResults(query, options.Category, newResults, ttlConf.Time)
			if uerr != nil {
				// Error in updating cache is not returned, just logged
				log.Error().
					Caller().
					Err(uerr).
					Str("queryAnon", anonymize.String(query)).
					Str("queryHash", anonymize.HashToSHA256B64(query)).
					Msg("Error replacing old results while updating database")
			}
		}
	}
}
