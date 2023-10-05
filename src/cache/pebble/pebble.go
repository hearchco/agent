package pebble

import (
	"fmt"

	"github.com/cockroachdb/pebble"
	"github.com/fxamacker/cbor/v2"
	"github.com/rs/zerolog/log"
	"github.com/tminaorg/brzaguza/src/cache"
)

type DB struct {
	pdb *pebble.DB
}

func New(path string) *DB {
	pdb, err := pebble.Open(fmt.Sprintf("%v/database", path), &pebble.Options{})
	if err != nil {
		log.Panic().Msgf("Error opening pebble: %v (path: %v)", err, path)
	} else {
		log.Info().Msgf("Successful connection to pebble (path: %v)", path)
	}
	return &DB{pdb: pdb}
}

func (db *DB) Close() {
	if err := db.pdb.Close(); err != nil {
		log.Panic().Msgf("Error closing pebble: %v", err)
	} else {
		log.Debug().Msg("Successfully disconnected from pebble")
	}
}

func (db *DB) Set(k string, v cache.Value) {
	if val, err := cbor.Marshal(v); err != nil {
		log.Error().Msgf("Error marshaling value: %v", err)
	} else if err := db.pdb.Set([]byte(k), val, pebble.Sync); err != nil { // should be set to NoSync when router gets graceful shutdown
		log.Panic().Msgf("Error setting KV to pebble: %v", err)
	}
}

func (db *DB) Get(k string, o cache.Value) {
	v, c, err := db.pdb.Get([]byte(k))
	val := []byte(v) // copy data before closing, casting needed for unmarshal
	if err == pebble.ErrNotFound {
		log.Trace().Msgf("Found no value in pebble for key (%v): %v", k, err)
	} else if err != nil {
		log.Panic().Msgf("Error getting value from pebble for key (%v): %v", k, err)
	} else if err := c.Close(); err != nil {
		log.Panic().Msgf("Error closing connection to pebble for key (%v): %v", k, err)
	} else if err := cbor.Unmarshal(val, o); err != nil {
		log.Error().Msgf("Failed unmarshaling value from pebble for key (%v): %v", k, err)
	}
}
