package nocache

import (
	"time"

	"github.com/hearchco/hearchco/src/cache"
)

type DB struct{}

func New() (*DB, error) { return nil, nil }

func (db *DB) Close() {}

func (db *DB) Set(k string, v cache.Value, ttl ...time.Duration) error { return nil }

func (db *DB) Get(k string, o cache.Value, hashed ...bool) error { return nil }

func (db *DB) GetTTL(k string, hashed ...bool) (time.Duration, error) { return 0, nil }
