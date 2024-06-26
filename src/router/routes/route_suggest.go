package routes

import (
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/hearchco/agent/src/config"
	"github.com/hearchco/agent/src/search"
	"github.com/hearchco/agent/src/search/engines/options"
	"github.com/hearchco/agent/src/search/result"
	"github.com/hearchco/agent/src/search/result/rank"
)

func routeSuggest(w http.ResponseWriter, r *http.Request, ver string, catConf config.Category) error {
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

	// Rank the suggestions.
	var rankedSugs rank.Suggestions = slices.Clone(scrapedSugs)
	rankedSugs.Rank(catConf.Ranking)

	// Check if the response should be in API format or normal JSON format.
	if strings.Contains(r.Header.Get("Accept"), "application/x-suggestions+json") {
		// Convert the suggestions to slice of strings.
		stringSugs := result.ConvertSuggestionsToOutput(rankedSugs)

		// If writing response failes, return the error.
		return writeResponseSuggestions(w, http.StatusOK, query, stringSugs)
	} else {
		// Create the response.
		res := SuggestionsResponse{
			responseBase{
				ver,
				time.Since(startTime).Milliseconds(),
			},
			rankedSugs,
		}

		// If writing response failes, return the error.
		return writeResponseJSON(w, http.StatusOK, res)
	}
}
