package cache

import (
	"time"
)

type Driver interface {
	Close()
	Set(k string, v any, ttl ...time.Duration) error
	Get(k string, o any) error
	GetTTL(k string) (time.Duration, error)
}
