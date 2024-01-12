package router

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/rs/zerolog/log"

	"github.com/hearchco/hearchco/src/bucket/result"
	"github.com/hearchco/hearchco/src/cache"
	"github.com/hearchco/hearchco/src/category"
	"github.com/hearchco/hearchco/src/config"
	"github.com/hearchco/hearchco/src/engines"
	"github.com/hearchco/hearchco/src/search"
)

func Search(c *gin.Context, conf *config.Config, db cache.DB) error {
	var query, pages, deepSearch, locale, categ, useragent, safesearch, mobile string
	var ccateg category.Name

	switch c.Request.Method {
	case "", "GET":
		{
			query = c.Query("q")
			pages = c.DefaultQuery("pages", "1")
			deepSearch = c.DefaultQuery("deep", "false")
			locale = c.DefaultQuery("locale", config.DefaultLocale)
			categ = c.DefaultQuery("category", "")
			useragent = c.DefaultQuery("useragent", "")
			safesearch = c.DefaultQuery("safesearch", "false")
			mobile = c.DefaultQuery("mobile", "false")
		}
	case "POST":
		{
			query = c.PostForm("q")
			pages = c.DefaultPostForm("pages", "1")
			deepSearch = c.DefaultPostForm("deep", "false")
			locale = c.DefaultPostForm("locale", config.DefaultLocale)
			categ = c.DefaultPostForm("category", "")
			useragent = c.DefaultPostForm("useragent", "")
			safesearch = c.DefaultPostForm("safesearch", "false")
			mobile = c.DefaultPostForm("mobile", "false")
		}
	}

	if query == "" {
		c.String(http.StatusOK, "")
	} else {
		maxPages, pageserr := strconv.Atoi(pages)
		if pageserr != nil {
			c.String(http.StatusUnprocessableEntity, fmt.Sprintf("Cannot convert pages value (\"%v\") to int", pages))
			return fmt.Errorf("router.Search(): cannot convert pages value \"%v\" to int: %w", pages, pageserr)
		}

		visitPages, deeperr := strconv.ParseBool(deepSearch)
		if deeperr != nil {
			c.String(http.StatusUnprocessableEntity, fmt.Sprintf("Cannot convert deep value (\"%v\") to bool", deepSearch))
			return fmt.Errorf("router.Search(): cannot convert deep value \"%v\" to int: %w", deepSearch, deeperr)
		}

		if lerr := engines.ValidateLocale(locale); lerr != nil {
			c.String(http.StatusUnprocessableEntity, fmt.Sprintf("Invalid locale value (\"%v\"), should be of the form \"en_US\"", locale))
			return fmt.Errorf("router.Search(): invalid locale value \"%v\": %w", locale, lerr)
		}

		ccateg = category.SafeFromString(categ)
		if ccateg == category.UNDEFINED {
			c.String(http.StatusUnprocessableEntity, fmt.Sprintf("Invalid category value (\"%v\")", categ))
			return fmt.Errorf("router.Search(): invalid category value \"%v\"", categ)
		}

		safeSearchB, safeerr := strconv.ParseBool(safesearch)
		if safeerr != nil {
			c.String(http.StatusUnprocessableEntity, fmt.Sprintf("Cannot convert safesearch value (\"%v\") to bool", safesearch))
			return fmt.Errorf("router.Search(): cannot convert safesearch value \"%v\" to bool: %w", safesearch, safeerr)
		}

		isMobile, mobileerr := strconv.ParseBool(mobile)
		if mobileerr != nil {
			c.String(http.StatusUnprocessableEntity, fmt.Sprintf("Cannot convert mobile value (\"%v\") to bool", mobile))
			return fmt.Errorf("router.Search(): cannot convert mobile value \"%v\" to bool: %w", mobile, mobileerr)
		}

		options := engines.Options{
			MaxPages:   maxPages,
			VisitPages: visitPages,
			Category:   ccateg,
			UserAgent:  useragent,
			Locale:     locale,
			SafeSearch: safeSearchB,
			Mobile:     isMobile,
		}

		var results []result.Result
		var foundInDB bool
		gerr := db.Get(query, &results)
		if gerr != nil {
			// Error in reading cache is not returned, just logged
			log.Error().
				Err(gerr).
				Str("Query", query).
				Msg("router.Search(): failed accessing cache")
		} else if results != nil {
			foundInDB = true
		} else {
			foundInDB = false
		}

		if foundInDB {
			log.Debug().
				Str("Query", query).
				Msg("Found results in cache")
		} else {
			log.Debug().Msg("Nothing found in cache, doing a clean search")

			results = search.PerformSearch(query, options, conf)
		}

		resultsShort := result.Shorten(results)
		if resultsJson, err := json.Marshal(resultsShort); err != nil {
			c.String(http.StatusInternalServerError, "")
			return fmt.Errorf("router.Search(): failed marshalling results: %v\n with error: %w", resultsShort, err)
		} else {
			c.String(http.StatusOK, string(resultsJson))
		}

		if !foundInDB {
			serr := db.Set(query, results)
			if serr != nil {
				// Error in updating cache is not returned, just logged
				log.Error().
					Err(serr).
					Str("Query", query).
					Msg("router.Search(): error updating database with search results")
			}
		}
	}
	return nil
}
