package moreurls

import (
	"fmt"
	"net/url"

	"github.com/rs/zerolog/log"
)

// Param struct, a simple key/value string pair.
type Param struct {
	key   string
	value string
}

// Constructs a new param with provided key and value.
func NewParam(k, v string) Param {
	return Param{k, v}
}

// Returns the key of the param.
func (p Param) Key() string {
	return p.key
}

// Return the value of the param.
func (p Param) Value() string {
	return p.value
}

// Sets the value of the param.
func (p *Param) SetValue(v string) {
	p.value = v
}

// Returns a copy of the param.
func (p Param) Copy() Param {
	return Param{p.key, p.value}
}

// Returns raw param in format "foo=bar".
func (p Param) String() string {
	if p.key == "" || p.value == "" {
		log.Panic().
			Str("key", p.key).
			Str("value", p.value).
			Msg("Empty key or value in parameter")
		// ^PANIC - Assert proper KV pair in Param.
	}
	return fmt.Sprintf("%s=%s", p.key, p.value)
}

// Returns URL encoded param in format "foo=bar".
func (p Param) QueryEscape() string {
	if p.key == "" || p.value == "" {
		log.Panic().
			Str("key", p.key).
			Str("value", p.value).
			Msg("Empty key or value in parameter")
		// ^PANIC - Assert proper KV pair in Param.
	}
	return fmt.Sprintf("%s=%s", url.QueryEscape(p.key), url.QueryEscape(p.value))
}
