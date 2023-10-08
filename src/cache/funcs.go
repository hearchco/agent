package cache

import (
	"time"

	"github.com/rs/zerolog/log"
	"github.com/tminaorg/brzaguza/src/bucket/result"
)

func Save(db DB, query string, results []result.Result) {
	log.Debug().Msg("Caching...")
	cacheTimer := time.Now()
	db.Set(query, results)
	log.Debug().Msgf("Cached results in %vns", time.Since(cacheTimer).Nanoseconds())
}
