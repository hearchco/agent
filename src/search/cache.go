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
	query string, cat category.Name, options engines.Options, db cache.DB,
	ttlConf config.TTL, settings map[engines.Name]config.Settings, categories map[category.Name]config.Category,
	results []result.Result, foundInDB bool,
	salt string,
) {
	if !foundInDB {
		log.Debug().
			Str("queryAnon", anonymize.String(query)).
			Str("queryHash", anonymize.HashToSHA256B64(query)).
			Msg("Caching results...")

		var err error
		if cat == category.IMAGES {
			err = db.SetImageResults(query, results)
		} else {
			err = db.SetResults(query, results)
		}

		if err != nil {
			log.Error().
				Err(err).
				Str("queryAnon", anonymize.String(query)).
				Str("queryHash", anonymize.HashToSHA256B64(query)).
				Msg("cli.Run(): error updating database with search results")
		}
	} else {
		log.Debug().
			Str("queryAnon", anonymize.String(query)).
			Str("queryHash", anonymize.HashToSHA256B64(query)).
			Msg("Checking if results need to be updated")

		age, err := db.GetAge(query)
		if err != nil {
			log.Error().
				Err(err).
				Str("queryAnon", anonymize.String(query)).
				Str("queryHash", anonymize.HashToSHA256B64(query)).
				Msg("cli.Run(): error getting TTL from database")
		}

		if age > ttlConf.RefreshTime {
			log.Info().
				Str("queryAnon", anonymize.String(query)).
				Str("queryHash", anonymize.HashToSHA256B64(query)).
				Msg("Updating results...")

			newResults := PerformSearch(query, options, settings, categories, salt)
			var err error

			// TODO: make this first delete old results and then add new ones
			// or update the old results with new data (?)
			if cat == category.IMAGES {
				err = db.SetImageResults(query, newResults)
			} else {
				err = db.SetResults(query, newResults)
			}

			if err != nil {
				// Error in updating cache is not returned, just logged
				log.Error().
					Err(err).
					Str("queryAnon", anonymize.String(query)).
					Str("queryHash", anonymize.HashToSHA256B64(query)).
					Msg("cli.Run(): error replacing old results while updating database")
			}
		}
	}
}
