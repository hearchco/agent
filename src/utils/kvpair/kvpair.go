package kvpair

import (
	"net/url"

	"github.com/rs/zerolog/log"
)

// KVPair struct, a simple key/value string pair.
type KVPair struct {
	key   string
	value string
}

// Constructs a new KVPair with provided key and value.
func NewKVPair(k, v string) KVPair {
	kv := KVPair{k, v}
	kv.assert()
	return kv
}

// Private assert function to ensure key and value are not empty.
// Panics if either key or value are empty.
func (kv KVPair) assert() {
	if kv.key == "" || kv.value == "" {
		log.Panic().
			Str("key", kv.key).
			Str("value", kv.value).
			Msg("Empty key or value in KVPair")
		// ^PANIC - Assert proper values in KVPair.
	}
}

// Returns the key.
func (kv KVPair) Key() string {
	kv.assert()
	return kv.key
}

// Returns the value.
func (kv KVPair) Value() string {
	kv.assert()
	return kv.value
}

// Sets the value.
func (kv *KVPair) SetValue(v string) {
	kv.assert()
	kv.value = v
	kv.assert()
}

// Returns a copy of the KVPair.
func (kv KVPair) Copy() KVPair {
	kv.assert()
	return NewKVPair(kv.key, kv.value)
}

// Returns raw KVPair in format "foo=bar".
func (kv KVPair) String() string {
	kv.assert()
	return kv.key + "=" + kv.value
}

// Returns URL encoded KVPair in format "foo=bar".
// Calls url.QueryEscape on both key and value.
func (kv KVPair) QueryEscape() string {
	kv.assert()
	return url.QueryEscape(kv.key) + "=" + url.QueryEscape(kv.value)
}
