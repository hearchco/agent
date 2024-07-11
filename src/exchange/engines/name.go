package engines

import (
	"fmt"
	"strings"
)

type Name int

const (
	UNDEFINED Name = iota
	CURRENCYAPI
	EXCHANGERATEAPI
	FRANKFURTER
)

func (n Name) String() string {
	switch n {
	case CURRENCYAPI:
		return "CurrencyAPI"
	case EXCHANGERATEAPI:
		return "ExchangeRateAPI"
	case FRANKFURTER:
		return "Frankfurter"
	default:
		return "Undefined"
	}
}

func (n Name) ToLower() string {
	return strings.ToLower(n.String())
}

func NameString(s string) (Name, error) {
	switch strings.ToLower(s) {
	case CURRENCYAPI.ToLower():
		return CURRENCYAPI, nil
	case EXCHANGERATEAPI.ToLower():
		return EXCHANGERATEAPI, nil
	case FRANKFURTER.ToLower():
		return FRANKFURTER, nil
	default:
		return UNDEFINED, fmt.Errorf("%s does not belong to Name values", s)
	}
}
