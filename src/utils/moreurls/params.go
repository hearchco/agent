package moreurls

import (
	"github.com/rs/zerolog/log"

	"github.com/hearchco/agent/src/utils/morestrings"
)

// Params struct, slice of Param struct.
type Params struct {
	params []Param
}

// Constructs a new slice of params with provided keys and values.
// Input should be in pairs: "key1, value1, key2, value2, key3, value3, ..."
// Number of elements must be even, otherwise this function panics.
func NewParams(elem ...string) Params {
	if len(elem)%2 != 0 {
		log.Panic().
			Strs("elem", elem).
			Msg("Odd number of elements for KV pairs")
		// ^PANIC - Assert even number of elements.
	}

	// Extract keys and values from elements.
	length := len(elem) / 2
	keys := make([]string, 0, length)
	values := make([]string, 0, length)
	for i, e := range elem {
		if i%2 == 0 {
			keys = append(keys, e)
		} else {
			values = append(values, e)
		}
	}

	// Create parameters slice.
	p := make([]Param, 0, length)
	for i := range length {
		p = append(p, Param{keys[i], values[i]})
	}

	return Params{p}
}

// Returns the value of the first occurence with the provided key.
// If not found, returns empty string and false.
func (p Params) Get(k string) (string, bool) {
	for _, param := range p.params {
		if param.Key() != k {
			continue
		}

		return param.Value(), true
	}

	return "", false
}

// Sets the value to the first occurence of the provided key.
// If not found, appends new Param KV pair and returns false.
func (p *Params) Set(k, v string) bool {
	for i, param := range p.params {
		if param.key != k {
			continue
		}

		p.params[i].SetValue(v)
		return true
	}

	p.params = append(p.params, Param{k, v})
	return false
}

// Returns a copy of the slice of params.
func (p Params) Copy() Params {
	n := make([]Param, 0, len(p.params))
	for _, param := range p.params {
		n = append(n, param.Copy())
	}
	return Params{n}
}

// Returns raw params in format "foo=bar&baz=woo".
func (p Params) String() string {
	paramsArray := make([]string, 0, len(p.params))
	for _, param := range p.params {
		paramsArray = append(paramsArray, param.String())
	}
	return morestrings.JoinNonEmpty("", "&", paramsArray...)
}

// Returns URL encoded params in format "foo=bar&baz=woo".
func (p Params) QueryEscape() string {
	paramsArray := make([]string, 0, len(p.params))
	for _, param := range p.params {
		paramsArray = append(paramsArray, param.QueryEscape())
	}
	return morestrings.JoinNonEmpty("", "&", paramsArray...)
}
