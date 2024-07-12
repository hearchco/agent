package currencyapi

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/hearchco/agent/src/exchange/currency"
)

// Rates field is named the same as base currency.
func (e Exchange) extractRates(resp string, base currency.Currency) (map[string]float64, error) {
	pattern := `"` + base.Lower() + `":\s*{[^}]*}`
	regexp := regexp.MustCompile(pattern)
	match := regexp.FindString(resp)
	if match == "" {
		return nil, fmt.Errorf("could not find JSON field for base currency %s", base)
	}

	// Remove `"<base_currency>":`` from the match
	jsonRates := strings.TrimSpace((match[len(base.Lower())+3:]))

	var rates map[string]float64
	if err := json.Unmarshal([]byte(jsonRates), &rates); err != nil {
		return nil, fmt.Errorf("could not unmarshal JSON field for base currency %s: %w", base, err)
	}

	return rates, nil
}
