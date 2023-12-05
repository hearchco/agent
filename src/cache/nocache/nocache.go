package nocache

import (
	"github.com/hearchco/hearchco/src/cache"
)

type DB struct{}

func New() *DB { return nil }

func (db *DB) Close() {}

func (db *DB) Set(k string, v cache.Value) {}

func (db *DB) Get(k string, o cache.Value) {}
