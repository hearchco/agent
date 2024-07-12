package currency

import (
	"fmt"
	"slices"
	"strings"

	"github.com/rs/zerolog/log"
)

// Format: ISO 4217 (3-letter code) e.g. CHF, EUR, GBP, USD.
type Currency string

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

func ConvertBase(curr string) Currency {
	// Hardcoded to ensure all APIs include these currencies and therefore work as expected.
	supportedBaseCurrencies := [...]string{"CHF", "EUR", "GBP", "USD"}

	upperCurr := strings.ToUpper(curr)
	if !slices.Contains(supportedBaseCurrencies[:], upperCurr) {
		log.Panic().
			Str("currency", upperCurr).
			Msg("unsupported base currency")
		// ^PANIC
	}

	return Currency(upperCurr)
}
