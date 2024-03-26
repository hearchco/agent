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
)

// returns response body, header and error
func Search(w http.ResponseWriter, r *http.Request, db cache.DB, ttlConf config.TTL, settings map[engines.Name]config.Settings, categories map[category.Name]config.Category, salt string) error {
	err := r.ParseForm()
	if err != nil {
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

	params := r.Form

	query := getParamOrDefault(params, "q")
	pagesStartS := getParamOrDefault(params, "start", "1")
	pagesMaxS := getParamOrDefault(params, "pages", "1")
	visitPagesS := getParamOrDefault(params, "deep", "false")
	locale := getParamOrDefault(params, "locale", config.DefaultLocale)
	categoryS := getParamOrDefault(params, "category", "")
	userAgent := getParamOrDefault(params, "useragent", "")
	safeSearchS := getParamOrDefault(params, "safesearch", "false")
	mobileS := getParamOrDefault(params, "mobile", "false")

	if query == "" || strings.TrimSpace(query) == "" || category.FromQuery(query) != "" {
		// return "[]" JSON when the query is empty or contains only category name
		return writeResponseJSON(w, http.StatusOK, []struct{}{})
	}

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

	visitPages, err := strconv.ParseBool(visitPagesS)
	if err != nil {
		// user error
		return writeResponseJSON(w, http.StatusUnprocessableEntity, ErrorResponse{
			Message: "cannot convert deep value to bool",
			Value:   fmt.Sprintf("%v", err),
		})
	}

	err = engines.ValidateLocale(locale)
	if err != nil {
		// user error
		return writeResponseJSON(w, http.StatusBadRequest, ErrorResponse{
			Message: "invalid locale value",
			Value:   fmt.Sprintf("%v", err),
		})
	}

	categoryName := category.SafeFromString(categoryS)
	if categoryName == category.UNDEFINED {
		// user error
		return writeResponseJSON(w, http.StatusBadRequest, ErrorResponse{
			Message: "invalid category value",
			Value:   fmt.Sprintf("%v", category.UNDEFINED),
		})
	}

	safeSearch, err := strconv.ParseBool(safeSearchS)
	if err != nil {
		// user error
		return writeResponseJSON(w, http.StatusUnprocessableEntity, ErrorResponse{
			Message: "cannot convert safesearch value to bool",
			Value:   fmt.Sprintf("%v", err),
		})
	}

	mobile, err := strconv.ParseBool(mobileS)
	if err != nil {
		// user error
		return writeResponseJSON(w, http.StatusUnprocessableEntity, ErrorResponse{
			Message: "cannot convert mobile value to bool",
			Value:   fmt.Sprintf("%v", err),
		})
	}

	options := engines.Options{
		Pages: engines.Pages{
			Start: pagesStart,
			Max:   pagesMax,
		},
		VisitPages: visitPages,
		Category:   categoryName,
		UserAgent:  userAgent,
		Locale:     locale,
		SafeSearch: safeSearch,
		Mobile:     mobile,
	}

	// search for results in db and web, afterwards return JSON
	results, foundInDB := search.Search(query, options, db, settings, categories, salt)

	// send response as soon as possible
	err = writeResponseJSON(w, http.StatusOK, results)

	// don't return immediately, we want to cache results and update them if necessary
	search.CacheAndUpdateResults(query, options, db, ttlConf, settings, categories, results, foundInDB, salt)

	// if writing response failed, return the error
	return err
}
