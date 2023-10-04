package cache

import (
	"time"

	"github.com/rs/zerolog/log"
	"github.com/tminaorg/brzaguza/src/bucket/result"
)

func Save(db DB, query string, results *result.Results) {
	log.Debug().Msg("Caching...")
	cacheTiming := time.Now()
	db.Set(query, results)
	log.Debug().Msgf("Cached results in %vms", time.Since(cacheTiming).Milliseconds())
}
