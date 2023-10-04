package cache

type DB interface {
	Close()
	Set(k string, v Value)
	Get(k string, o Value)
}

type Value interface{}
