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
		log.Fatal().Msgf("Error opening pebble: %v (path: %v)", err, pebblePath)
	} else {
		log.Info().Msgf("Successfully opened pebble (path: %v)", pebblePath)
	}

	return &DB{pdb: pdb}
}

func (db *DB) Close() {
	if err := db.pdb.Close(); err != nil {
		log.Fatal().Msgf("Error closing pebble: %v", err)
	} else {
		log.Debug().Msg("Successfully closed pebble")
	}
}

func (db *DB) Set(k string, v cache.Value) {
	log.Debug().Msg("Caching...")
	cacheTimer := time.Now()

	if val, err := cbor.Marshal(v); err != nil {
		log.Error().Msgf("Error marshaling value: %v", err)
	} else if err := db.pdb.Set([]byte(k), val, pebble.NoSync); err != nil {
		log.Fatal().Msgf("Error setting KV to pebble: %v", err)
	} else {
		cacheTimeSince := time.Since(cacheTimer)
		log.Debug().Msgf("Cached results in %vms (%vns)", cacheTimeSince.Milliseconds(), cacheTimeSince.Nanoseconds())
	}
}

func (db *DB) Get(k string, o cache.Value) {
	v, c, err := db.pdb.Get([]byte(k))
	val := []byte(v) // copy data before closing, casting needed for unmarshal

	if err == pebble.ErrNotFound {
		log.Trace().Msgf("Found no value in pebble for key (%v): %v", k, err)
	} else if err != nil {
		log.Fatal().Msgf("Error getting value from pebble for key (%v): %v", k, err)
	} else if err := c.Close(); err != nil {
		log.Fatal().Msgf("Error closing io to pebble for key (%v): %v", k, err)
	} else if err := cbor.Unmarshal(val, o); err != nil {
		log.Error().Msgf("Failed unmarshaling value from pebble for key (%v): %v", k, err)
	}
}
