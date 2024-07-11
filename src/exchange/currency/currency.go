package currency

import (
	"fmt"
	"strings"
)

// Format: ISO 4217 (3-letter code) e.g. USD, EUR, GBP.
type Currency string

// TODO: Make base currency configurable.
const Base Currency = "EUR"

func (c Currency) String() string {
	return string(c)
}

func (c Currency) Lower() string {
	return strings.ToLower(c.String())
}

func Convert(curr string) (Currency, error) {
	if len(curr) != 3 {
		return "", fmt.Errorf("currency code must be 3 characters long")
	}

	upperCurr := strings.ToUpper(curr)
	return Currency(upperCurr), nil
}
