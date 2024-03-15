package router

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/hearchco/hearchco/src/cache"
	"github.com/hearchco/hearchco/src/config"
	"github.com/hearchco/hearchco/src/gotypelimits"
	"github.com/hearchco/hearchco/src/search"
	"github.com/hearchco/hearchco/src/search/category"
	"github.com/hearchco/hearchco/src/search/engines"
)

func Search(c *fiber.Ctx, db cache.DB, ttlConf config.TTL, settings map[engines.Name]config.Settings, categories map[category.Name]config.Category) error {
	var query, pagesStartS, pagesMaxS, visitPagesS, locale, categoryS, userAgent, safeSearchS, mobileS string

	switch c.Method() {
	case "", "GET":
		{
			query = c.Query("q")
			pagesStartS = c.Query("start", "1")
			pagesMaxS = c.Query("pages", "1")
			visitPagesS = c.Query("deep", "false")
			locale = c.Query("locale", config.DefaultLocale)
			categoryS = c.Query("category", "")
			userAgent = c.Query("useragent", "")
			safeSearchS = c.Query("safesearch", "false")
			mobileS = c.Query("mobile", "false")
		}
	case "POST":
		{
			query = c.FormValue("q")
			pagesStartS = c.FormValue("start", "1")
			pagesMaxS = c.FormValue("pages", "1")
			visitPagesS = c.FormValue("deep", "false")
			locale = c.FormValue("locale", config.DefaultLocale)
			categoryS = c.FormValue("category", "")
			userAgent = c.FormValue("useragent", "")
			safeSearchS = c.FormValue("safesearch", "false")
			mobileS = c.FormValue("mobile", "false")
		}
	}

	// TODO: implement more cases when query is useless to process
	if query == "" {
		// return empty array of objects
		return c.Status(fiber.StatusBadRequest).JSON([]struct{}{})
	}

	pagesMax, err := strconv.Atoi(pagesMaxS)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(ErrorResponse{
			Message: "Cannot convert pages value to int",
			Value:   pagesMaxS,
		})
	}
	// TODO: make upper limit configurable
	pagesMaxUpperLimit := 10
	if pagesMax < 1 || pagesMax > pagesMaxUpperLimit {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Message: fmt.Sprintf("Pages value must be at least 1 and at most %v", pagesMaxUpperLimit),
			Value:   pagesMaxS,
		})
	}

	pagesStart, err := strconv.Atoi(pagesStartS)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(ErrorResponse{
			Message: "Cannot convert start value to int",
			Value:   pagesStartS,
		})
	}
	// make sure that pagesStart can be safely added to pagesMax
	if pagesStart < 1 || pagesStart > gotypelimits.MaxInt-pagesMaxUpperLimit {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Message: "Start value must be at least 1",
			Value:   pagesStartS,
		})
	} else {
		// since it's >=1, we decrement it to match the 0-based index
		pagesStart -= 1
	}

	visitPages, err := strconv.ParseBool(visitPagesS)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(ErrorResponse{
			Message: "Cannot convert deep value to bool",
			Value:   visitPagesS,
		})
	}

	err = engines.ValidateLocale(locale)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Message: "Invalid locale value, should be of the form \"en_US\"",
			Value:   locale,
		})
	}

	categoryName := category.SafeFromString(categoryS)
	if categoryName == category.UNDEFINED {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Message: "Invalid category value",
			Value:   categoryS,
		})
	}

	safeSearch, err := strconv.ParseBool(safeSearchS)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(ErrorResponse{
			Message: "Cannot convert safesearch value to bool",
			Value:   safeSearchS,
		})
	}

	mobile, err := strconv.ParseBool(mobileS)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(ErrorResponse{
			Message: "Cannot convert mobile value to bool",
			Value:   mobileS,
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
	results, foundInDB := search.Search(query, options, db, settings, categories)
	err = c.JSON(results)

	// run even if response errored out (dropped connection?)
	// check if results need to be cached/updated and if so do it
	search.CacheAndUpdateResults(query, options, db, ttlConf, settings, categories, results, foundInDB)

	return err
}
