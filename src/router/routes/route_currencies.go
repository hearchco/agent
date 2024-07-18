package routes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/hearchco/agent/src/cache"
	"github.com/hearchco/agent/src/config"
	"github.com/hearchco/agent/src/exchange"
	"github.com/rs/zerolog/log"
)

func routeCurrencies(w http.ResponseWriter, ver string, conf config.Exchange, db cache.DB, ttl time.Duration) error {
	// Capture start time.
	startTime := time.Now()

	// Get the cached currencies.
	currencies, err := db.GetCurrencies(conf.BaseCurrency, conf.Engines)
	if err != nil {
		log.Error().
			Err(err).
			Str("base", conf.BaseCurrency.String()).
			Str("engines", fmt.Sprintf("%v", conf.Engines)).
			Msg("Error while getting currencies from cache")
	}

	// Create the exchange.
	var exch exchange.Exchange
	if currencies == nil {
		// Fetch the currencies from the enabled engines.
		exch = exchange.NewExchange(conf)
		// Cache the currencies if any have been fetched.
		if len(exch.Currencies()) > 0 {
			err := db.SetCurrencies(conf.BaseCurrency, conf.Engines, exch.Currencies(), ttl)
			if err != nil {
				log.Error().
					Err(err).
					Str("base", conf.BaseCurrency.String()).
					Str("engines", fmt.Sprintf("%v", conf.Engines)).
					Msg("Error while setting currencies in cache")
			}
		}
	} else {
		// Use the cached currencies.
		exch = exchange.NewExchange(conf, currencies)
	}

	return writeResponseJSON(w, http.StatusOK, CurrenciesResponse{
		responseBase{
			ver,
			time.Since(startTime).Milliseconds(),
		},
		conf.BaseCurrency,
		exch.Currencies(),
	})
}
