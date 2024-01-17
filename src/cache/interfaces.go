package cache

type DB interface {
	Close()
	Set(k string, v Value) error
	Get(k string, o Value, hashed ...bool) error
}

type Value interface{}
