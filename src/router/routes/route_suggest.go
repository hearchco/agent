package routes

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hearchco/agent/src/search"
	"github.com/hearchco/agent/src/search/engines/options"
	"github.com/hearchco/agent/src/search/result"
)

func routeSuggest(w http.ResponseWriter, r *http.Request) error {
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

	// Query is required.
	query := strings.TrimSpace(getParamOrDefault(r.Form, "q"))
	if query == "" {
		// User error.
		return writeResponseJSON(w, http.StatusBadRequest, ErrorResponse{
			Message: "query cannot be empty or whitespace",
			Value:   "empty query",
		})
	}

	localeS := getParamOrDefault(r.Form, "locale", options.LocaleDefault.String())
	locale, err := options.StringToLocale(localeS)
	if err != nil {
		// User error.
		return writeResponseJSON(w, http.StatusBadRequest, ErrorResponse{
			Message: "invalid locale value",
			Value:   fmt.Sprintf("%v", err),
		})
	}

	// Search for suggestions.
	// TODO: Make timeout configurable.
	scrapedSugs, err := search.Suggest(query, locale, 1*time.Second)
	if err != nil {
		// Server error.
		werr := writeResponseJSON(w, http.StatusInternalServerError, ErrorResponse{
			Message: "failed to search",
			Value:   fmt.Sprintf("%v", err),
		})
		if werr != nil {
			return fmt.Errorf("%w: %w", werr, err)
		}
		return err
	}

	// Rank the suggestions.
	// rankedSugs := rank.Rank(scrapedSugs, Ranking)

	outputSugs := result.ConvertSuggestionsToOutput(scrapedSugs)

	// Create the response.
	res := [...]any{
		query,
		outputSugs,
	}

	// If writing response failes, return the error.
	return writeResponseJSON(w, http.StatusOK, res)
}
