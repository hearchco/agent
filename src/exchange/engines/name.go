package engines

import (
	"strings"
)

type Name int

//go:generate enumer -type=Name -json -text -sql
//go:generate go run github.com/hearchco/agent/generate/exchanger -type=Name -packagename exchange -output ../engine_exchanger.go
const (
	UNDEFINED Name = iota
	CURRENCYAPI
	EXCHANGERATEAPI
	FRANKFURTER
)

// Returns engine names without UNDEFINED.
func Names() []Name {
	return _NameValues[1:]
}

func (n Name) ToLower() string {
	return strings.ToLower(n.String())
}
