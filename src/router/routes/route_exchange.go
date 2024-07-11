package routes

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hearchco/agent/src/exchange"
	"github.com/hearchco/agent/src/exchange/currency"
)

func routeExchange(w http.ResponseWriter, r *http.Request, ver string) error {
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

	// TODO: Make base currency configurable.
	const base currency.Currency = "EUR"
	convertedAmount, err := exchange.Exchange(base, from, to, amount)
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

	return writeResponseJSON(w, http.StatusOK, ExchangeResponse{
		responseBase{
			ver,
			time.Since(startTime).Milliseconds(),
		},
		base,
		from,
		to,
		amount,
		convertedAmount,
	})
}
