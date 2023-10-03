package cache

type DB interface {
	Close()
	Set(k string, v interface{})
	Get(k string) string
}
