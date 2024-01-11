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
		log.Error().Err(err).Msgf("pebble.New(): error opening pebble at path: %v", pebblePath)
	} else {
		log.Info().Msgf("Successfully opened pebble (path: %v)", pebblePath)
	}

	return &DB{pdb: pdb}
}

func (db *DB) Close() {
	if err := db.pdb.Close(); err != nil {
		log.Error().Err(err).Msg("pebble.Close(): error closing pebble")
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
		return fmt.Errorf("pebble.Set(): error setting KV to pebble: %w", err)
	} else {
		cacheTimeSince := time.Since(cacheTimer)
		log.Trace().Msgf("Cached results in %vms (%vns)", cacheTimeSince.Milliseconds(), cacheTimeSince.Nanoseconds())
	}
	return nil
}

func (db *DB) Get(k string, o cache.Value) error {
	v, c, err := db.pdb.Get([]byte(k))
	val := []byte(v) // copy data before closing, casting needed for unmarshal

	if err == pebble.ErrNotFound {
		log.Trace().Msgf("Found no value in pebble for key: \"%v\"", k)
	} else if err != nil {
		return fmt.Errorf("pebble.Get(): error getting value from pebble for key %v: %w", k, err)
	} else if err := c.Close(); err != nil {
		return fmt.Errorf("pebble.Get(): error closing io to pebble for key %v: %w", k, err)
	} else if err := cbor.Unmarshal(val, o); err != nil {
		return fmt.Errorf("pebble.Get(): failed unmarshaling value from pebble for key %v: %w", k, err)
	}
	return nil
}
