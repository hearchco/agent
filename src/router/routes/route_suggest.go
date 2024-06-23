package routes

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/hearchco/agent/src/config"
	"github.com/hearchco/agent/src/search"
	"github.com/hearchco/agent/src/search/engines/options"
	"github.com/hearchco/agent/src/search/result"
)

func routeSuggest(w http.ResponseWriter, r *http.Request, catConf config.Category) error {
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
	scrapedSugs, err := search.Suggest(query, locale, catConf)
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

	// TODO: Rank the suggestions.
	// rankedSugs := rank.Rank(scrapedSugs, Ranking)

	// Convert the suggestions to output format.
	outputSugs := result.ConvertSuggestionsToOutput(scrapedSugs)

	// Check if the response should be in API format or normal JSON format.
	api := strings.Contains(r.Header.Get("Accept"), "application/x-suggestions+json")

	// If writing response failes, return the error.
	return writeResponseSuggestions(w, http.StatusOK, query, outputSugs, api)
}
