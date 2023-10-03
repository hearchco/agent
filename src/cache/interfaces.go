package cache

type DB interface {
	Close()
	Set(k string, v Value)
	Get(k string) []byte
}

type Value interface {
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(data []byte) error
}
