package router

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/hearchco/hearchco/src/cache"
	"github.com/hearchco/hearchco/src/config"
	"github.com/hearchco/hearchco/src/gotypelimits"
	"github.com/hearchco/hearchco/src/search"
	"github.com/hearchco/hearchco/src/search/category"
	"github.com/hearchco/hearchco/src/search/engines"
	"github.com/hearchco/hearchco/src/search/result"
)

// returns response body, header and error
func Search(w http.ResponseWriter, r *http.Request, db cache.DB, ttlConf config.TTL, settings map[engines.Name]config.Settings, categories map[category.Name]config.Category, salt string) error {
	// parse form data (including query params)
	if err := r.ParseForm(); err != nil {
		// server error
		werr := writeResponseJSON(w, http.StatusInternalServerError, ErrorResponse{
			Message: "failed to parse form",
			Value:   fmt.Sprintf("%v", err),
		})
		if werr != nil {
			return fmt.Errorf("%w: %w", werr, err)
		}
		return err
	}

	query := strings.TrimSpace(getParamOrDefault(r.Form, "q")) // query is required
	if query == "" {
		// user error
		return writeResponseJSON(w, http.StatusBadRequest, ErrorResponse{
			Message: "query cannot be empty or whitespace",
			Value:   "empty query",
		})
	}

	visitPagesS := getParamOrDefault(r.Form, "deep", "false")
	visitPages, err := strconv.ParseBool(visitPagesS)
	if err != nil {
		// user error
		return writeResponseJSON(w, http.StatusUnprocessableEntity, ErrorResponse{
			Message: "cannot convert deep value to bool",
			Value:   fmt.Sprintf("%v", err),
		})
	}

	safeSearchS := getParamOrDefault(r.Form, "safesearch", "false")
	safeSearch, err := strconv.ParseBool(safeSearchS)
	if err != nil {
		// user error
		return writeResponseJSON(w, http.StatusUnprocessableEntity, ErrorResponse{
			Message: "cannot convert safesearch value to bool",
			Value:   fmt.Sprintf("%v", err),
		})
	}

	pagesMaxS := getParamOrDefault(r.Form, "pages", "1")
	pagesMax, err := strconv.Atoi(pagesMaxS)
	if err != nil {
		// user error
		return writeResponseJSON(w, http.StatusUnprocessableEntity, ErrorResponse{
			Message: "cannot convert pages value to int",
			Value:   fmt.Sprintf("%v", err),
		})
	}
	// TODO: make upper limit configurable
	pagesMaxUpperLimit := 10
	if pagesMax < 1 || pagesMax > pagesMaxUpperLimit {
		// user error
		return writeResponseJSON(w, http.StatusBadRequest, ErrorResponse{
			Message: fmt.Sprintf("pages value must be at least 1 and at most %v", pagesMaxUpperLimit),
			Value:   "out of range",
		})
	}

	pagesStartS := getParamOrDefault(r.Form, "start", "1")
	pagesStart, err := strconv.Atoi(pagesStartS)
	if err != nil {
		// user error
		return writeResponseJSON(w, http.StatusUnprocessableEntity, ErrorResponse{
			Message: "cannot convert start value to int",
			Value:   fmt.Sprintf("%v", err),
		})
	}
	// make sure that pagesStart can be safely added to pagesMax
	if pagesStart < 1 || pagesStart > gotypelimits.MaxInt-pagesMaxUpperLimit {
		// user error
		return writeResponseJSON(w, http.StatusBadRequest, ErrorResponse{
			Message: fmt.Sprintf("start value must be at least 1 and at most %v", gotypelimits.MaxInt-pagesMaxUpperLimit),
			Value:   "out of range",
		})
	} else {
		// since it's >=1, we decrement it to match the 0-based index
		pagesStart -= 1
	}

	locale := getParamOrDefault(r.Form, "locale", config.DefaultLocale)
	err = engines.ValidateLocale(locale)
	if err != nil {
		// user error
		return writeResponseJSON(w, http.StatusBadRequest, ErrorResponse{
			Message: "invalid locale value",
			Value:   fmt.Sprintf("%v", err),
		})
	}

	categoryS := getParamOrDefault(r.Form, "category", category.GENERAL.String())
	categoryName, err := category.FromString(categoryS)
	if err != nil {
		// user error
		return writeResponseJSON(w, http.StatusBadRequest, ErrorResponse{
			Message: "invalid category value",
			Value:   fmt.Sprintf("%v", categoryName),
		})
	}

	// all of these have default values set and are validated beforehand
	options := engines.Options{
		VisitPages: visitPages,
		SafeSearch: safeSearch,
		Pages: engines.Pages{
			Start: pagesStart,
			Max:   pagesMax,
		},
		Locale:   locale,
		Category: categoryName,
	}

	// search for results
	results, foundInDB := search.Search(query, options, db, categories[options.Category], settings, salt)

	// send response as soon as possible
	if categoryName == category.IMAGES {
		resultsOutput := result.ConvertToImageOutput(results)
		err = writeResponseJSON(w, http.StatusOK, resultsOutput)
	} else {
		resultsOutput := result.ConvertToGeneralOutput(results)
		err = writeResponseJSON(w, http.StatusOK, resultsOutput)
	}

	// TODO: this doesn't work on AWS Lambda because the response is already sent (which terminates the process)
	// don't return immediately, we want to cache results and update them if necessary
	search.CacheAndUpdateResults(query, options, db, ttlConf, categories[options.Category], settings, results, foundInDB, salt)

	// if writing response failed, return the error
	return err
}
