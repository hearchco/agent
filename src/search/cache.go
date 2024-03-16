package search

import (
	"github.com/hearchco/hearchco/src/anonymize"
	"github.com/hearchco/hearchco/src/cache"
	"github.com/hearchco/hearchco/src/config"
	"github.com/hearchco/hearchco/src/search/category"
	"github.com/hearchco/hearchco/src/search/engines"
	"github.com/hearchco/hearchco/src/search/result"
	"github.com/rs/zerolog/log"
)

func CacheAndUpdateResults(
	query string, options engines.Options, db cache.DB,
	ttlConf config.TTL, settings map[engines.Name]config.Settings, categories map[category.Name]config.Category,
	results []result.Result, foundInDB bool,
	salt string,
) {
	if !foundInDB {
		log.Debug().
			Str("queryAnon", anonymize.String(query)).
			Str("queryHash", anonymize.HashToSHA256B64(query)).
			Msg("Caching results...")
		serr := db.Set(query, results, ttlConf.Time)
		if serr != nil {
			log.Error().
				Err(serr).
				Str("queryAnon", anonymize.String(query)).
				Str("queryHash", anonymize.HashToSHA256B64(query)).
				Msg("cli.Run(): error updating database with search results")
		}
	} else {
		log.Debug().
			Str("queryAnon", anonymize.String(query)).
			Str("queryHash", anonymize.HashToSHA256B64(query)).
			Msg("Checking if results need to be updated")
		ttl, terr := db.GetTTL(query)
		if terr != nil {
			log.Error().
				Err(terr).
				Str("queryAnon", anonymize.String(query)).
				Str("queryHash", anonymize.HashToSHA256B64(query)).
				Msg("cli.Run(): error getting TTL from database")
		} else if ttl < ttlConf.RefreshTime {
			log.Info().
				Str("queryAnon", anonymize.String(query)).
				Str("queryHash", anonymize.HashToSHA256B64(query)).
				Msg("Updating results...")
			newResults := PerformSearch(query, options, settings, categories, salt)
			uerr := db.Set(query, newResults, ttlConf.Time)
			if uerr != nil {
				// Error in updating cache is not returned, just logged
				log.Error().
					Err(uerr).
					Str("queryAnon", anonymize.String(query)).
					Str("queryHash", anonymize.HashToSHA256B64(query)).
					Msg("cli.Run(): error replacing old results while updating database")
			}
		}
	}
}
