package cache

import "time"

type DB interface {
	Close()
	Set(k string, v interface{}, ttl ...time.Duration) error
	Get(k string, o interface{}, hashed ...bool) error
	GetTTL(k string, hashed ...bool) (time.Duration, error)
}
