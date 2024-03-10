package router

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/hearchco/hearchco/src/cache"
	"github.com/hearchco/hearchco/src/config"
	"github.com/hearchco/hearchco/src/gotypelimits"
	"github.com/hearchco/hearchco/src/search"
	"github.com/hearchco/hearchco/src/search/category"
	"github.com/hearchco/hearchco/src/search/engines"
	"github.com/hearchco/hearchco/src/search/result"
)

func Search(c *gin.Context, db cache.DB, ttlConf config.TTL, settings map[engines.Name]config.Settings, categories map[category.Name]config.Category) error {
	var query, pagesStartS, pagesMaxS, visitPagesS, locale, categoryS, userAgent, safeSearchS, mobileS string

	switch c.Request.Method {
	case "", "GET":
		{
			query = c.Query("q")
			pagesStartS = c.DefaultQuery("start", "1")
			pagesMaxS = c.DefaultQuery("pages", "1")
			visitPagesS = c.DefaultQuery("deep", "false")
			locale = c.DefaultQuery("locale", config.DefaultLocale)
			categoryS = c.DefaultQuery("category", "")
			userAgent = c.DefaultQuery("useragent", "")
			safeSearchS = c.DefaultQuery("safesearch", "false")
			mobileS = c.DefaultQuery("mobile", "false")
		}
	case "POST":
		{
			query = c.PostForm("q")
			pagesStartS = c.DefaultPostForm("start", "1")
			pagesMaxS = c.DefaultPostForm("pages", "1")
			visitPagesS = c.DefaultPostForm("deep", "false")
			locale = c.DefaultPostForm("locale", config.DefaultLocale)
			categoryS = c.DefaultPostForm("category", "")
			userAgent = c.DefaultPostForm("useragent", "")
			safeSearchS = c.DefaultPostForm("safesearch", "false")
			mobileS = c.DefaultPostForm("mobile", "false")
		}
	}

	// TODO: implement more cases when query is useless to process
	if query == "" {
		// return empty array of objects
		c.JSON(http.StatusOK, []struct{}{})
		return nil
	}

	pagesMax, err := strconv.Atoi(pagesMaxS)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
			Message: "Cannot convert pages value to int",
			Value:   pagesMaxS,
		})
		return fmt.Errorf("router.Search(): cannot convert pages value %q to int: %w", pagesMaxS, err)
	}
	// TODO: make upper limit configurable
	pagesMaxUpperLimit := 10
	if pagesMax < 1 || pagesMax > pagesMaxUpperLimit {
		c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
			Message: fmt.Sprintf("Pages value must be at least 1 and at most %v", pagesMaxUpperLimit),
			Value:   pagesMaxS,
		})
		return fmt.Errorf("router.Search(): pages value %q must be at least 1 and at most %v", pagesMaxS, pagesMaxUpperLimit)
	}

	pagesStart, err := strconv.Atoi(pagesStartS)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
			Message: "Cannot convert start value to int",
			Value:   pagesStartS,
		})
		return fmt.Errorf("router.Search(): cannot convert start value %q to int: %w", pagesStartS, err)
	}
	// make sure that pagesStart can be safely added to pagesMax
	if pagesStart < 1 || pagesStart > gotypelimits.MaxInt-pagesMaxUpperLimit {
		c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
			Message: "Start value must be at least 1",
			Value:   pagesStartS,
		})
		return fmt.Errorf("router.Search(): start value %q must be at least 1", pagesStartS)
	} else {
		// since it's >=1, we decrement it to match the 0-based index
		pagesStart -= 1
	}

	visitPages, err := strconv.ParseBool(visitPagesS)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
			Message: "Cannot convert deep value to bool",
			Value:   visitPagesS,
		})
		return fmt.Errorf("router.Search(): cannot convert deep value %q to int: %w", visitPagesS, err)
	}

	err = engines.ValidateLocale(locale)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
			Message: "Invalid locale value, should be of the form \"en_US\"",
			Value:   locale,
		})
		return fmt.Errorf("router.Search(): invalid locale value %q: %w", locale, err)
	}

	categoryName := category.SafeFromString(categoryS)
	if categoryName == category.UNDEFINED {
		c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
			Message: "Invalid category value",
			Value:   categoryS,
		})
		return fmt.Errorf("router.Search(): invalid category value %q", categoryS)
	}

	safeSearch, err := strconv.ParseBool(safeSearchS)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
			Message: "Cannot convert safesearch value to bool",
			Value:   safeSearchS,
		})
		return fmt.Errorf("router.Search(): cannot convert safesearch value %q to bool: %w", safeSearchS, err)
	}

	mobile, err := strconv.ParseBool(mobileS)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
			Message: "Cannot convert mobile value to bool",
			Value:   mobileS,
		})
		return fmt.Errorf("router.Search(): cannot convert mobile value %q to bool: %w", mobileS, err)
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

	results, foundInDB := search.Search(query, options, db, settings, categories)
	c.JSON(http.StatusOK, result.Shorten(results))

	search.CacheAndUpdateResults(query, options, db, ttlConf, settings, categories, results, foundInDB)
	return nil
}
