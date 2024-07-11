package exchange

import (
	"fmt"
	"sync"

	"github.com/hearchco/agent/src/exchange/currency"
	"github.com/hearchco/agent/src/exchange/engines"
	"github.com/rs/zerolog/log"
)

func Exchange(base currency.Currency, from currency.Currency, to currency.Currency, amount float64) (float64, error) {
	// TODO: Load currency map from cache if available for the given engines (they should be part of the PK).

	enabledEngines := []engines.Name{engines.CURRENCYAPI, engines.EXCHANGERATEAPI, engines.FRANKFURTER}
	exchangers := exchangerArray()
	currencyArrayMap := currency.NewCurrencyMap()

	var wg sync.WaitGroup
	wg.Add(enginerLen - 1) // -1 because of UNDEFINED
	for _, eng := range enabledEngines {
		exch := exchangers[eng]
		go func() {
			defer wg.Done()
			currs, err := exch.Exchange(base)
			if err != nil {
				log.Error().
					Err(err).
					Str("engine", eng.String()).
					Msg("Error while exchanging")
				return
			}
			currencyArrayMap.Append(currs)
		}()
	}
	wg.Wait()

	// TODO: Cache the currency map.
	// Extract the averaged currency map.
	currencyMap := currencyArrayMap.Extract()

	// Check if FROM and TO currencies are supported.
	if _, ok := currencyMap[from]; !ok {
		return -1, fmt.Errorf("unsupported FROM currency: %s", from)
	}
	if _, ok := currencyMap[to]; !ok {
		return -1, fmt.Errorf("unsupported TO currency: %s", to)
	}

	// Convert the amount in FROM currency to base currency.
	basedAmount := amount / currencyMap[from]

	// Convert the amount in base currency to TO currency.
	convertedAmount := basedAmount * currencyMap[to]

	return convertedAmount, nil
}
