package search

import (
	"github.com/hearchco/hearchco/src/anonymize"
	"github.com/hearchco/hearchco/src/bucket/result"
	"github.com/hearchco/hearchco/src/cache"
	"github.com/hearchco/hearchco/src/config"
	"github.com/hearchco/hearchco/src/engines"
	"github.com/rs/zerolog/log"
)

func CacheAndUpdateResults(query string, options engines.Options, conf *config.Config, db cache.DB, results []result.Result, foundInDB bool) {
	if !foundInDB {
		log.Debug().
			Str("queryAnon", anonymize.String(query)).
			Str("queryHash", anonymize.HashToSHA256B64(query)).
			Msg("Caching results...")
		serr := db.Set(query, results, conf.Server.Cache.TTL.Time)
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
		} else if ttl < conf.Server.Cache.TTL.RefreshTime {
			log.Info().
				Str("queryAnon", anonymize.String(query)).
				Str("queryHash", anonymize.HashToSHA256B64(query)).
				Msg("Updating results...")
			newResults := performSearch(query, options, conf)
			uerr := db.Set(query, newResults, conf.Server.Cache.TTL.Time)
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
