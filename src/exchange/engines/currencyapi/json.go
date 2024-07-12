package currencyapi

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/hearchco/agent/src/exchange/currency"
)

// Rates field is named the same as base currency.
func extractRatesFromResp(resp string, base currency.Currency) (map[string]float64, error) {
	pattern := `"` + base.Lower() + `":\s*{[^}]*}`
	regexp := regexp.MustCompile(pattern)
	match := regexp.FindString(resp)
	if match == "" {
		return nil, fmt.Errorf("could not find JSON field for base currency %s", base)
	}

	var rates map[string]float64
	if err := json.Unmarshal([]byte(match), &rates); err != nil {
		return nil, fmt.Errorf("could not unmarshal JSON field for base currency %s: %w", base, err)
	}

	return rates, nil
}
