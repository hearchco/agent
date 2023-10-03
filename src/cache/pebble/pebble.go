package pebble

import (
	"fmt"

	"github.com/cockroachdb/pebble"
	"github.com/rs/zerolog/log"
)

type DB struct {
	pdb *pebble.DB
}

func New(path string) *DB {
	pdb, err := pebble.Open(fmt.Sprintf("%v/database", path), &pebble.Options{})
	if err != nil {
		log.Panic().Msgf("error opening pebble: %v (path: %v)", err, path)
	} else {
		log.Info().Msgf("successful connection to pebble (path: %v)", path)
	}
	return &DB{pdb: pdb}
}

func (db *DB) Close() {
	if err := db.pdb.Close(); err != nil {
		log.Panic().Msgf("error closing pebble: %v", err)
	}
}

func (db *DB) Set(k string, v interface{}) {
	if err := db.pdb.Set([]byte(k), []byte(fmt.Sprint(v)), pebble.Sync); err != nil {
		log.Trace().Msgf("key = %v, value = %v", k, v)
		log.Panic().Msgf("error setting KV to pebble: %v", err)
	} else {
		log.Trace().Msgf("success: key = %v, value = %v", k, v)
	}
}

func (db *DB) Get(k string) string {
	v, c, err := db.pdb.Get([]byte(k))
	val := string(v) // copy data before closing, casting needed for json.Unmarshal()

	if err != nil {
		log.Trace().Msgf("error: key = %v, value = %v", k, val)
		log.Panic().Msgf("error getting value from pebble for key (%v): %v", k, err)
	}
	if err := c.Close(); err != nil {
		log.Trace().Msgf("error: key = %v, value = %v", k, val)
		log.Panic().Msgf("error closing connection to pebble for key (%v): %v", k, err)
	}

	log.Trace().Msgf("success: key = %v, value = %v", k, val)
	return val
}
