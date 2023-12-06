package pebble

import (
	"fmt"
	"path"
	"time"

	"github.com/cockroachdb/pebble"
	"github.com/fxamacker/cbor/v2"
	"github.com/hearchco/hearchco/src/cache"
	"github.com/rs/zerolog/log"
)

type DB struct {
	pdb *pebble.DB
}

func New(dataDirPath string) *DB {
	pebblePath := path.Join(dataDirPath, "database")
	pdb, err := pebble.Open(pebblePath, &pebble.Options{})

	if err != nil {
		log.Fatal().Err(err).Msgf("pebble.New(): error opening pebble at path: %v", pebblePath)
		// ^FATAL
	} else {
		log.Info().Msgf("Successfully opened pebble (path: %v)", pebblePath)
	}

	return &DB{pdb: pdb}
}

func (db *DB) Close() {
	if err := db.pdb.Close(); err != nil {
		log.Fatal().Err(err).Msg("pebble.Close(): error closing pebble")
		// ^FATAL
	} else {
		log.Debug().Msg("Successfully closed pebble")
	}
}

func (db *DB) Set(k string, v cache.Value) error {
	log.Debug().Msg("Caching...")
	cacheTimer := time.Now()

	if val, err := cbor.Marshal(v); err != nil {
		return fmt.Errorf("pebble.Set(): error marshaling value: %w", err)
	} else if err := db.pdb.Set([]byte(k), val, pebble.NoSync); err != nil {
		log.Fatal().Err(err).Msg("pebble.Set(): error setting KV to pebble")
		// ^FATAL
	} else {
		cacheTimeSince := time.Since(cacheTimer)
		log.Debug().Msgf("Cached results in %vms (%vns)", cacheTimeSince.Milliseconds(), cacheTimeSince.Nanoseconds())
	}
	return nil
}

func (db *DB) Get(k string, o cache.Value) error {
	v, c, err := db.pdb.Get([]byte(k))
	val := []byte(v) // copy data before closing, casting needed for unmarshal

	if err == pebble.ErrNotFound {
		log.Trace().Msgf("Found no value in pebble for key: \"%v\"", k)
	} else if err != nil {
		log.Fatal().Err(err).Msgf("pebble.Get(): error getting value from pebble for key %v", k)
		// ^FATAL
	} else if err := c.Close(); err != nil {
		log.Fatal().Err(err).Msgf("pebble.Get(): error closing io to pebble for key %v", k)
		// ^FATAL
	} else if err := cbor.Unmarshal(val, o); err != nil {
		return fmt.Errorf("pebble.Get(): failed unmarshaling value from pebble for key %v: %w", k, err)
	}
	return nil
}
