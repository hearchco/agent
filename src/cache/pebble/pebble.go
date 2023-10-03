package pebble

import (
	"github.com/cockroachdb/pebble"
	"github.com/rs/zerolog/log"
)

type DB struct {
	pdb *pebble.DB
}

func New(path string) *DB {
	pdb, err := pebble.Open(path, &pebble.Options{})
	if err != nil {
		log.Panic().Msgf("error opening pebble: %v (path: %v)", err, path)
	}
	return &DB{pdb: pdb}
}

func (db *DB) Close() {
	if err := db.pdb.Close(); err != nil {
		log.Panic().Msgf("error closing pebble: %v", err)
	}
}

func (db *DB) Set(k, v []byte) {
	if err := db.pdb.Set(k, v, pebble.Sync); err != nil {
		log.Trace().Msgf("key = %v\nvalue = %v", k, v)
		log.Panic().Msgf("error setting KV to pebble: %v", err)
	}
}

func (db *DB) Get(k []byte) []byte {
	if v, c, err := db.pdb.Get(k); err != nil {
		log.Trace().Msgf("error: key = %v\nvalue = %v", k, v)
		log.Panic().Msgf("error getting value from pebble for key (%v): %v", k, err)
	} else if err := c.Close(); err != nil {
		log.Trace().Msgf("error: key = %v\nvalue = %v", k, v)
		log.Panic().Msgf("error closing connection to pebble for key (%v): %v", k, err)
	} else {
		log.Trace().Msgf("success: key = %v\nvalue = %v", k, v)
		return v
	}
	return nil
}
