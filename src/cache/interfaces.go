package cache

type DB interface {
	Close()
	Set(k string, v Value) error
	Get(k string, o Value) error
}

type Value interface{}
