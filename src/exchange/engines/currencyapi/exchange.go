package currencyapi

import (
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
	dataRates, err := extractRatesFromResp(string(body), base)
	if err != nil {
		return nil, fmt.Errorf("failed to extract rates from response: %w", err)
	}

	// Check if no rates were found.
	if len(dataRates) == 0 {
		return nil, fmt.Errorf("no rates found for %s", base)
	}

	// Convert the rates to proper currency types with their rates.
	rates := make(currency.Currencies, len(dataRates))
	for currS, rate := range dataRates {
		curr, err := currency.Convert(currS)
		if err != nil {
			// Non-ISO currencies are expected from this engine.
			log.Trace().
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
