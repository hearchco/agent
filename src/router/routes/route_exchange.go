package routes

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hearchco/agent/src/cache"
	"github.com/hearchco/agent/src/config"
	"github.com/hearchco/agent/src/exchange"
	"github.com/hearchco/agent/src/exchange/currency"
	"github.com/rs/zerolog/log"
)

func routeExchange(w http.ResponseWriter, r *http.Request, ver string, conf config.Exchange, db cache.DB, ttl time.Duration) error {
	// Capture start time.
	startTime := time.Now()

	// Parse form data (including query params).
	if err := r.ParseForm(); err != nil {
		// Server error.
		werr := writeResponseJSON(w, http.StatusInternalServerError, ErrorResponse{
			Message: "failed to parse form",
			Value:   fmt.Sprintf("%v", err),
		})
		if werr != nil {
			return fmt.Errorf("%w: %w", werr, err)
		}
		return err
	}

	// FROM is required.
	fromS := strings.TrimSpace(getParamOrDefault(r.Form, "from"))
	if fromS == "" {
		// User error.
		return writeResponseJSON(w, http.StatusBadRequest, ErrorResponse{
			Message: "from cannot be empty or whitespace",
			Value:   "empty from",
		})
	}

	// Parse FROM currency.
	from, err := currency.Convert(fromS)
	if err != nil {
		// User error.
		return writeResponseJSON(w, http.StatusBadRequest, ErrorResponse{
			Message: "invalid from currency",
			Value:   fmt.Sprintf("%v", err),
		})
	}

	// TO is required.
	toS := strings.TrimSpace(getParamOrDefault(r.Form, "to"))
	if toS == "" {
		// User error.
		return writeResponseJSON(w, http.StatusBadRequest, ErrorResponse{
			Message: "to cannot be empty or whitespace",
			Value:   "empty to",
		})
	}

	// Parse TO currency.
	to, err := currency.Convert(toS)
	if err != nil {
		// User error.
		return writeResponseJSON(w, http.StatusBadRequest, ErrorResponse{
			Message: "invalid to currency",
			Value:   fmt.Sprintf("%v", err),
		})
	}

	// AMOUNT is required.
	amountS := strings.TrimSpace(getParamOrDefault(r.Form, "amount"))
	if amountS == "" {
		// User error.
		return writeResponseJSON(w, http.StatusBadRequest, ErrorResponse{
			Message: "amount cannot be empty or whitespace",
			Value:   "empty amount",
		})
	}

	// Parse amount.
	amount, err := strconv.ParseFloat(amountS, 64)
	if err != nil {
		// User error.
		return writeResponseJSON(w, http.StatusBadRequest, ErrorResponse{
			Message: "invalid amount value",
			Value:   fmt.Sprintf("%v", err),
		})
	}

	// Get the cached currencies.
	currencies, err := db.GetCurrencies(currency.Base, conf.Engines)
	if err != nil {
		log.Error().
			Err(err).
			Str("base", currency.Base.String()).
			Str("engines", fmt.Sprintf("%v", conf.Engines)).
			Msg("Error while getting currencies from cache")
	}

	// Create the exchange.
	var exch exchange.Exchange
	if currencies == nil {
		// Fetch the currencies from the enabled engines.
		exch = exchange.NewExchange(currency.Base, conf)
		// Cache the currencies if any have been fetched.
		if len(exch.Currencies()) > 0 {
			err := db.SetCurrencies(currency.Base, conf.Engines, exch.Currencies(), ttl)
			if err != nil {
				log.Error().
					Err(err).
					Str("base", currency.Base.String()).
					Str("engines", fmt.Sprintf("%v", conf.Engines)).
					Msg("Error while setting currencies in cache")
			}
		}
	} else {
		// Use the cached currencies.
		exch = exchange.NewExchange(currency.Base, conf, currencies)
	}

	// Check if FROM and TO are supported currencies.
	if !exch.SupportsCurrency(from) {
		// User error.
		return writeResponseJSON(w, http.StatusBadRequest, ErrorResponse{
			Message: "unsupported from currency",
			Value:   from.String(),
		})
	}
	if !exch.SupportsCurrency(to) {
		// User error.
		return writeResponseJSON(w, http.StatusBadRequest, ErrorResponse{
			Message: "unsupported to currency",
			Value:   to.String(),
		})
	}

	// Convert the amount.
	convAmount := exch.Convert(from, to, amount)

	return writeResponseJSON(w, http.StatusOK, ExchangeResponse{
		responseBase{
			ver,
			time.Since(startTime).Milliseconds(),
		},
		currency.Base,
		from,
		to,
		amount,
		convAmount,
	})
}
