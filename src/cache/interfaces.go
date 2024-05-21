package cache

import "time"

type DB interface {
	Close()
	Set(k string, v interface{}, ttl ...time.Duration) error
	Get(k string, o interface{}) error
	GetTTL(k string) (time.Duration, error)
}
