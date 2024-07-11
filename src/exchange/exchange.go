package exchange

import (
	"sync"

	"github.com/hearchco/agent/src/exchange/currency"
	"github.com/hearchco/agent/src/exchange/engines"
	"github.com/rs/zerolog/log"
)

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

func (e Exchange) SupportsCurrency(curr currency.Currency) bool {
	_, ok := e.currencies[curr]
	return ok
}

func (e Exchange) Convert(from currency.Currency, to currency.Currency, amount float64) float64 {
	// Check if FROM and TO are supported currencies.
	if !e.SupportsCurrency(from) || !e.SupportsCurrency(to) {
		log.Panic().
			Str("from", from.String()).
			Str("to", to.String()).
			Msg("Unsupported currencies")
		// ^PANIC - This should never happen.
	}

	// Convert the amount in FROM currency to base currency.
	basedAmount := amount / e.currencies[from]

	// Convert the amount in base currency to TO currency.
	convertedAmount := basedAmount * e.currencies[to]

	return convertedAmount
}
