package currencyapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/hearchco/agent/src/exchange/currency"
)

func (e Exchange) Exchange(base currency.Currency) (currency.Currencies, error) {
	// Get data from the API.
	api := e.apiUrlWithBaseCurrency(base)
	resp, err := http.Get(api)
	if err != nil {
		return nil, fmt.Errorf("failed to get data from %s: %w", api, err)
	}

	// Read the response body.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Unmarshal the response.
	var data response
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Check if no rates were found.
	if len(data.Rates) == 0 {
		return nil, fmt.Errorf("no rates found for %s", base)
	}

	// Convert the rates to proper currency types with their rates.
	rates := make(currency.Currencies, len(data.Rates))
	for currS, rate := range data.Rates {
		curr, err := currency.Convert(currS)
		if err != nil {
			log.Error().
				Err(err).
				Str("currency", currS).
				Msg("failed to convert currency")
			continue
		}
		rates[curr] = rate
	}

	// Set the base currency rate to 1.
	rates[base] = 1

	return rates, nil
}
