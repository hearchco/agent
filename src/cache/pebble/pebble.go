package pebble

import (
	"path"
	"time"

	"github.com/cockroachdb/pebble"
	"github.com/fxamacker/cbor/v2"
	"github.com/rs/zerolog/log"
	"github.com/tminaorg/brzaguza/src/cache"
)

type DB struct {
	pdb *pebble.DB
}

func New(dataDirPath string) *DB {
	pebblePath := path.Join(dataDirPath, "database")
	pdb, err := pebble.Open(pebblePath, &pebble.Options{})

	if err != nil {
		log.Fatal().Err(err).Msgf("pebble.New(): error opening pebble at path: %v", pebblePath)
		return nil
	} else {
		log.Info().Msgf("Successfully opened pebble (path: %v)", pebblePath)
	}

	return &DB{pdb: pdb}
}

func (db *DB) Close() {
	if err := db.pdb.Close(); err != nil {
		log.Fatal().Err(err).Msg("pebble.Close(): error closing pebble")
		return
	} else {
		log.Debug().Msg("Successfully closed pebble")
	}
}

func (db *DB) Set(k string, v cache.Value) {
	log.Debug().Msg("Caching...")
	cacheTimer := time.Now()

	if val, err := cbor.Marshal(v); err != nil {
		log.Error().Err(err).Msg("pebble.Set(): error marshaling value")
	} else if err := db.pdb.Set([]byte(k), val, pebble.NoSync); err != nil {
		log.Fatal().Err(err).Msg("pebble.Set(): error setting KV to pebble")
		return
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
		log.Fatal().Err(err).Msgf("pebble.Get(): error getting value from pebble for key %v", k)
		return
	} else if err := c.Close(); err != nil {
		log.Fatal().Err(err).Msgf("pebble.Get(): error closing io to pebble for key %v", k)
		return
	} else if err := cbor.Unmarshal(val, o); err != nil {
		log.Error().Err(err).Msgf("pebble.Get(): failed unmarshaling value from pebble for key %v", k)
	}
}
