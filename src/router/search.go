package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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
		writeResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse form: %v", err))
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

	// TODO: implement more cases when query is useless to process
	if query == "" {
		// return empty array of objects
		res, err := json.Marshal([]struct{}{})
		if err != nil {
			// server error
			writeResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to marshal empty array: %v", err))
			return err
		}
		writeResponseJSON(w, http.StatusOK, res)
		return nil
	}

	pagesMax, err := strconv.Atoi(pagesMaxS)
	if err != nil {
		// user error
		writeResponse(w, http.StatusUnprocessableEntity, fmt.Sprintf("cannot convert pages value to int: %v", err))
		return nil
	}

	// TODO: make upper limit configurable
	pagesMaxUpperLimit := 10
	if pagesMax < 1 || pagesMax > pagesMaxUpperLimit {
		// user error
		writeResponse(w, http.StatusBadRequest, fmt.Sprintf("pages value must be at least 1 and at most %v", pagesMaxUpperLimit))
		return nil
	}

	pagesStart, err := strconv.Atoi(pagesStartS)
	if err != nil {
		// user error
		writeResponse(w, http.StatusUnprocessableEntity, fmt.Sprintf("cannot convert start value to int: %v", err))
		return nil
	}

	// make sure that pagesStart can be safely added to pagesMax
	if pagesStart < 1 || pagesStart > gotypelimits.MaxInt-pagesMaxUpperLimit {
		// user error
		writeResponse(w, http.StatusBadRequest, fmt.Sprintf("start value must be at least 1 and at most %v", gotypelimits.MaxInt-pagesMaxUpperLimit))
		return nil
	} else {
		// since it's >=1, we decrement it to match the 0-based index
		pagesStart -= 1
	}

	visitPages, err := strconv.ParseBool(visitPagesS)
	if err != nil {
		// user error
		writeResponse(w, http.StatusUnprocessableEntity, fmt.Sprintf("cannot convert deep value to bool: %v", err))
		return nil
	}

	err = engines.ValidateLocale(locale)
	if err != nil {
		// user error
		writeResponse(w, http.StatusBadRequest, fmt.Sprintf("invalid locale value: %v", err))
		return nil
	}

	categoryName := category.SafeFromString(categoryS)
	if categoryName == category.UNDEFINED {
		// user error
		writeResponse(w, http.StatusBadRequest, "invalid category value")
		return nil
	}

	safeSearch, err := strconv.ParseBool(safeSearchS)
	if err != nil {
		// user error
		writeResponse(w, http.StatusUnprocessableEntity, fmt.Sprintf("cannot convert safesearch value to bool: %v", err))
		return nil
	}

	mobile, err := strconv.ParseBool(mobileS)
	if err != nil {
		// user error
		writeResponse(w, http.StatusUnprocessableEntity, fmt.Sprintf("cannot convert mobile value to bool: %v", err))
		return nil
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
	results, foundInDB, cat := search.Search(query, options, db, settings, categories, salt)
	resultsJson, err := json.Marshal(results)
	if err != nil {
		// server error
		writeResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to marshal results: %v", err))
		return err
	}

	// send response as soon as possible
	writeResponseJSON(w, http.StatusOK, []byte(resultsJson))
	// don't return immediately, we want to cache results and update them if necessary
	search.CacheAndUpdateResults(query, cat, options, db, ttlConf, settings, categories, results, foundInDB, salt)

	return nil
}
