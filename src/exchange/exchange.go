package exchange

import (
	"fmt"
	"sync"

	"github.com/hearchco/agent/src/exchange/currency"
	"github.com/hearchco/agent/src/exchange/engines"
	"github.com/rs/zerolog/log"
)

// TODO: Test caching with private fields.
type Exchange struct {
	base       currency.Currency
	currencies currency.Currencies
}

func NewExchange(base currency.Currency, enabledEngines []engines.Name, currencies ...currency.Currencies) Exchange {
	// If currencies are provided, use them.
	if len(currencies) > 0 {
		return Exchange{
			base,
			currencies[0],
		}
	}

	// Otherwise, fetch the currencies from the enabled engines.
	exchangers := exchangerArray()
	currencyMap := currency.NewCurrencyMap()

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
			currencyMap.Append(currs)
		}()
	}
	wg.Wait()

	return Exchange{
		base,
		currencyMap.Extract(),
	}
}

func (e Exchange) Currencies() currency.Currencies {
	return e.currencies
}

func (e Exchange) Convert(from currency.Currency, to currency.Currency, amount float64) (float64, error) {
	// Check if FROM and TO currencies are supported.
	if _, ok := e.currencies[from]; !ok {
		return -1, fmt.Errorf("unsupported FROM currency: %s", from)
	}
	if _, ok := e.currencies[to]; !ok {
		return -1, fmt.Errorf("unsupported TO currency: %s", to)
	}

	// Convert the amount in FROM currency to base currency.
	basedAmount := amount / e.currencies[from]

	// Convert the amount in base currency to TO currency.
	convertedAmount := basedAmount * e.currencies[to]

	return convertedAmount, nil
}
