package router

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/hearchco/hearchco/src/cache"
	"github.com/hearchco/hearchco/src/config"
	"github.com/hearchco/hearchco/src/search"
	"github.com/hearchco/hearchco/src/search/category"
	"github.com/hearchco/hearchco/src/search/engines"
	"github.com/hearchco/hearchco/src/search/result"
)

func Search(c *gin.Context, db cache.DB, ttlConf config.TTL, settings map[engines.Name]config.Settings, categories map[category.Name]config.Category) error {
	var query, start, pages, deepSearch, locale, cats, useragent, safesearch, mobile string

	switch c.Request.Method {
	case "", "GET":
		{
			query = c.Query("q")
			start = c.DefaultQuery("start", "1")
			pages = c.DefaultQuery("pages", "1")
			deepSearch = c.DefaultQuery("deep", "false")
			locale = c.DefaultQuery("locale", config.DefaultLocale)
			cats = c.DefaultQuery("category", "")
			useragent = c.DefaultQuery("useragent", "")
			safesearch = c.DefaultQuery("safesearch", "false")
			mobile = c.DefaultQuery("mobile", "false")
		}
	case "POST":
		{
			query = c.PostForm("q")
			start = c.DefaultPostForm("start", "1")
			pages = c.DefaultPostForm("pages", "1")
			deepSearch = c.DefaultPostForm("deep", "false")
			locale = c.DefaultPostForm("locale", config.DefaultLocale)
			cats = c.DefaultPostForm("category", "")
			useragent = c.DefaultPostForm("useragent", "")
			safesearch = c.DefaultPostForm("safesearch", "false")
			mobile = c.DefaultPostForm("mobile", "false")
		}
	}

	if query == "" {
		// return empty array of objects
		c.JSON(http.StatusOK, []struct{}{})
		return nil
	}

	pagesStart, starterr := strconv.Atoi(start)
	if starterr != nil {
		c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
			Message: "Cannot convert start value to int",
			Value:   start,
		})
		return fmt.Errorf("router.Search(): cannot convert start value %q to int: %w", start, starterr)
	}
	if pagesStart < 1 {
		c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
			Message: "Start value must be at least 1",
			Value:   start,
		})
		return fmt.Errorf("router.Search(): start value %q must be at least 1", start)
	} else {
		// since it's >=1, we decrement it to match the 0-based index
		pagesStart -= 1
	}

	pagesMax, pageserr := strconv.Atoi(pages)
	if pageserr != nil {
		c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
			Message: "Cannot convert pages value to int",
			Value:   pages,
		})
		return fmt.Errorf("router.Search(): cannot convert pages value %q to int: %w", pages, pageserr)
	}
	if pagesMax < 1 {
		c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
			Message: "Pages value must be at least 1",
			Value:   pages,
		})
		return fmt.Errorf("router.Search(): pages value %q must be at least 1", pages)
	}

	visitPages, deeperr := strconv.ParseBool(deepSearch)
	if deeperr != nil {
		c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
			Message: "Cannot convert deep value to bool",
			Value:   deepSearch,
		})
		return fmt.Errorf("router.Search(): cannot convert deep value %q to int: %w", deepSearch, deeperr)
	}

	if lerr := engines.ValidateLocale(locale); lerr != nil {
		c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
			Message: "Invalid locale value, should be of the form \"en_US\"",
			Value:   locale,
		})
		return fmt.Errorf("router.Search(): invalid locale value %q: %w", locale, lerr)
	}

	cat := category.SafeFromString(cats)
	if cat == category.UNDEFINED {
		c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
			Message: "Invalid category value",
			Value:   cats,
		})
		return fmt.Errorf("router.Search(): invalid category value %q", cats)
	}

	safeSearchB, safeerr := strconv.ParseBool(safesearch)
	if safeerr != nil {
		c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
			Message: "Cannot convert safesearch value to bool",
			Value:   safesearch,
		})
		return fmt.Errorf("router.Search(): cannot convert safesearch value %q to bool: %w", safesearch, safeerr)
	}

	isMobile, mobileerr := strconv.ParseBool(mobile)
	if mobileerr != nil {
		c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
			Message: "Cannot convert mobile value to bool",
			Value:   mobile,
		})
		return fmt.Errorf("router.Search(): cannot convert mobile value %q to bool: %w", mobile, mobileerr)
	}

	options := engines.Options{
		Pages: engines.Pages{
			Start: pagesStart,
			Max:   pagesMax,
		},
		VisitPages: visitPages,
		Category:   cat,
		UserAgent:  useragent,
		Locale:     locale,
		SafeSearch: safeSearchB,
		Mobile:     isMobile,
	}

	results, foundInDB := search.Search(query, options, db, settings, categories)
	c.JSON(http.StatusOK, result.Shorten(results))

	search.CacheAndUpdateResults(query, options, db, ttlConf, settings, categories, results, foundInDB)
	return nil
}
