package cache

import "time"

type DB interface {
	Close()
	Set(k string, v Value, ttl ...time.Duration) error
	Get(k string, o Value, hashed ...bool) error
}

type Value interface{}
