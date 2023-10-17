package pebble

import (
	"time"

	"github.com/cockroachdb/pebble"
	"github.com/fxamacker/cbor/v2"
	"github.com/rs/zerolog/log"
	"github.com/tminaorg/brzaguza/src/cache"
)

type DB struct {
	pdb *pebble.DB
}

func New(path string) *DB {
	pebblePath := path + "/database"
	pdb, err := pebble.Open(pebblePath, &pebble.Options{})

	if err != nil {
		log.Fatal().Err(err).Msgf("Error opening pebble at path: %v", pebblePath)
	} else {
		log.Info().Msgf("Successfully opened pebble (path: %v)", pebblePath)
	}

	return &DB{pdb: pdb}
}

func (db *DB) Close() {
	if err := db.pdb.Close(); err != nil {
		log.Fatal().Err(err).Msg("Error closing pebble")
	} else {
		log.Debug().Msg("Successfully closed pebble")
	}
}

func (db *DB) Set(k string, v cache.Value) {
	log.Debug().Msg("Caching...")
	cacheTimer := time.Now()

	if val, err := cbor.Marshal(v); err != nil {
		log.Error().Err(err).Msg("Error marshaling value")
	} else if err := db.pdb.Set([]byte(k), val, pebble.NoSync); err != nil {
		log.Fatal().Err(err).Msg("Error setting KV to pebble")
	} else {
		cacheTimeSince := time.Since(cacheTimer)
		log.Debug().Msgf("Cached results in %vms (%vns)", cacheTimeSince.Milliseconds(), cacheTimeSince.Nanoseconds())
	}
}

func (db *DB) Get(k string, o cache.Value) {
	v, c, err := db.pdb.Get([]byte(k))
	val := []byte(v) // copy data before closing, casting needed for unmarshal

	if err == pebble.ErrNotFound {
		log.Trace().Msgf("Found no value in pebble for key %v", k)
	} else if err != nil {
		log.Fatal().Err(err).Msgf("Error getting value from pebble for key %v", k)
	} else if err := c.Close(); err != nil {
		log.Fatal().Err(err).Msgf("Error closing io to pebble for key %v", k)
	} else if err := cbor.Unmarshal(val, o); err != nil {
		log.Error().Err(err).Msgf("Failed unmarshaling value from pebble for key %v", k)
	}
}
