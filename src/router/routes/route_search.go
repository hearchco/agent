package routes

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/hearchco/agent/src/cache"
	"github.com/hearchco/agent/src/config"
	"github.com/hearchco/agent/src/search"
	"github.com/hearchco/agent/src/search/category"
	"github.com/hearchco/agent/src/search/engines/options"
	"github.com/hearchco/agent/src/search/result"
	"github.com/hearchco/agent/src/search/result/rank"
	"github.com/hearchco/agent/src/utils/anonymize"
	"github.com/hearchco/agent/src/utils/gotypelimits"
	"github.com/rs/zerolog/log"
)

func routeSearch(w http.ResponseWriter, r *http.Request, catsConf map[category.Name]config.Category, ttlConf config.TTL, db cache.DB, salt string) error {
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

	categoryS := getParamOrDefault(r.Form, "category", category.GENERAL.String())
	categoryName, err := category.FromString(categoryS)
	if err != nil {
		// User error.
		return writeResponseJSON(w, http.StatusBadRequest, ErrorResponse{
			Message: "invalid category value",
			Value:   fmt.Sprintf("%v", categoryName),
		})
	}

	pagesMaxS := getParamOrDefault(r.Form, "pages", "1")
	pagesMax, err := strconv.Atoi(pagesMaxS)
	if err != nil {
		// User error.
		return writeResponseJSON(w, http.StatusUnprocessableEntity, ErrorResponse{
			Message: "cannot convert pages value to int",
			Value:   fmt.Sprintf("%v", err),
		})
	}
	// TODO: Make upper limit configurable.
	pagesMaxUpperLimit := 10
	if pagesMax < 1 || pagesMax > pagesMaxUpperLimit {
		// User error.
		return writeResponseJSON(w, http.StatusBadRequest, ErrorResponse{
			Message: fmt.Sprintf("pages value must be at least 1 and at most %v", pagesMaxUpperLimit),
			Value:   "out of range",
		})
	}

	pagesStartS := getParamOrDefault(r.Form, "start", "1")
	pagesStart, err := strconv.Atoi(pagesStartS)
	if err != nil {
		// User error.
		return writeResponseJSON(w, http.StatusUnprocessableEntity, ErrorResponse{
			Message: "cannot convert start value to int",
			Value:   fmt.Sprintf("%v", err),
		})
	}
	// Make sure that pagesStart can be safely added to pagesMax.
	if pagesStart < 1 || pagesStart > gotypelimits.MaxInt-pagesMaxUpperLimit {
		// User error.
		return writeResponseJSON(w, http.StatusBadRequest, ErrorResponse{
			Message: fmt.Sprintf("start value must be at least 1 and at most %v", gotypelimits.MaxInt-pagesMaxUpperLimit),
			Value:   "out of range",
		})
	} else {
		// Since it's >=1, we decrement it to match the 0-based index.
		pagesStart -= 1
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

	safeSearchS := getParamOrDefault(r.Form, "safesearch", "false")
	safeSearch, err := strconv.ParseBool(safeSearchS)
	if err != nil {
		// User error.
		return writeResponseJSON(w, http.StatusUnprocessableEntity, ErrorResponse{
			Message: "cannot convert safesearch value to bool",
			Value:   fmt.Sprintf("%v", err),
		})
	}

	// All of these have default values set and validated.
	opts := options.Options{
		Pages: options.Pages{
			Start: pagesStart,
			Max:   pagesMax,
		},
		Locale:     locale,
		SafeSearch: safeSearch,
	}

	// Check cache for results.
	cachedRes, err := db.GetResults(query, categoryName, opts)
	if err != nil {
		log.Error().
			Err(err).
			Str("query", anonymize.String(query)).
			Str("category", categoryName.String()).
			Msg("Failed to get results from cache")
	} else if len(cachedRes) > 0 {
		log.Debug().
			Str("query", anonymize.String(query)).
			Str("category", categoryName.String()).
			Msg("Results found in cache")

		// Convert the results to include the hashes (output format).
		outpusRes := result.ConvertToOutput(cachedRes, salt)

		// If writing response failes, return the error.
		return writeResponseJSON(w, http.StatusOK, outpusRes)
	} else {
		log.Debug().
			Str("query", anonymize.String(query)).
			Str("category", categoryName.String()).
			Msg("No results found in cache")
	}

	// Search for results.
	scrapedRes, err := search.Search(query, categoryName, opts, catsConf[categoryName])
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

	// Rank the results.
	rankedRes := rank.Rank(scrapedRes, catsConf[categoryName].Ranking)

	// Store the results in cache.
	if err := db.SetResults(query, categoryName, opts, rankedRes, ttlConf.Time); err != nil {
		log.Error().
			Err(err).
			Str("query", anonymize.String(query)).
			Str("category", categoryName.String()).
			Msg("failed to set results in cache")
	}

	// Convert the results to include the hashes (output format).
	outpusRes := result.ConvertToOutput(rankedRes, salt)

	// If writing response failes, return the error.
	return writeResponseJSON(w, http.StatusOK, outpusRes)
}
