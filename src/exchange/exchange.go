package exchange

import (
	"context"
	"fmt"
	"sync"

	"github.com/hearchco/agent/src/config"
	"github.com/hearchco/agent/src/exchange/currency"
	"github.com/rs/zerolog/log"
)

type Exchange struct {
	base       currency.Currency
	currencies currency.Currencies
}

func NewExchange(base currency.Currency, conf config.Exchange, currencies ...currency.Currencies) Exchange {
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

	// Create context with HardTimeout.
	ctxHardTimeout, cancelHardTimeoutFunc := context.WithTimeout(context.Background(), conf.Timings.HardTimeout)
	defer cancelHardTimeoutFunc()

	// Create a WaitGroup for all engines.
	var wg sync.WaitGroup
	wg.Add(len(conf.Engines))

	// Create a context that cancels when the WaitGroup is done.
	exchangeCtx, cancelExchange := context.WithCancel(context.Background())
	defer cancelExchange()
	go func() {
		wg.Wait()
		cancelExchange()
	}()

	for _, eng := range conf.Engines {
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

	// Wait for either all engines to finish or the HardTimeout.
	select {
	case <-exchangeCtx.Done():
		log.Trace().
			Dur("timeout", conf.Timings.HardTimeout).
			Str("engines", fmt.Sprintf("%v", conf.Engines)).
			Msg("All engines finished")
	case <-ctxHardTimeout.Done():
		log.Trace().
			Dur("timeout", conf.Timings.HardTimeout).
			Str("engines", fmt.Sprintf("%v", conf.Engines)).
			Msg("HardTimeout reached")
	}

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
