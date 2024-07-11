package routes

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hearchco/agent/src/cache"
	"github.com/hearchco/agent/src/exchange"
	"github.com/hearchco/agent/src/exchange/currency"
	"github.com/hearchco/agent/src/exchange/engines"
	"github.com/rs/zerolog/log"
)

func routeExchange(w http.ResponseWriter, r *http.Request, ver string, db cache.DB, ttl time.Duration) error {
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

	// TODO: Make base currency and enabled engines configurable.
	const base currency.Currency = "EUR"
	enabledEngines := [...]engines.Name{engines.CURRENCYAPI, engines.EXCHANGERATEAPI, engines.FRANKFURTER}

	// Get the cached currencies.
	currencies, err := db.GetCurrencies(base, enabledEngines[:])
	if err != nil {
		log.Error().
			Err(err).
			Str("base", base.String()).
			Str("engines", fmt.Sprintf("%v", enabledEngines)).
			Msg("Error while getting currencies from cache")
	}

	// Create the exchange.
	var exch exchange.Exchange
	if currencies == nil {
		// Fetch the currencies from the enabled engines.
		exch = exchange.NewExchange(base, enabledEngines[:])
	} else {
		// Use the cached currencies.
		exch = exchange.NewExchange(base, enabledEngines[:], currencies)
	}

	// Convert the amount.
	convAmount, err := exch.Convert(from, to, amount)
	if err != nil {
		// Server error.
		werr := writeResponseJSON(w, http.StatusInternalServerError, ErrorResponse{
			Message: "failed to exchange",
			Value:   fmt.Sprintf("%v", err),
		})
		if werr != nil {
			return fmt.Errorf("%w: %w", werr, err)
		}
		return err
	}

	// Cache the currencies.
	if currencies == nil {
		err := db.SetCurrencies(base, enabledEngines[:], exch.Currencies(), ttl)
		if err != nil {
			log.Error().
				Err(err).
				Str("base", base.String()).
				Str("engines", fmt.Sprintf("%v", enabledEngines)).
				Msg("Error while setting currencies in cache")
		}
	}

	return writeResponseJSON(w, http.StatusOK, ExchangeResponse{
		responseBase{
			ver,
			time.Since(startTime).Milliseconds(),
		},
		base,
		from,
		to,
		amount,
		convAmount,
	})
}
