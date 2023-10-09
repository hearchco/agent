package router

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/rs/zerolog/log"

	"github.com/tminaorg/brzaguza/src/bucket/result"
	"github.com/tminaorg/brzaguza/src/cache"
	"github.com/tminaorg/brzaguza/src/config"
	"github.com/tminaorg/brzaguza/src/engines"
	"github.com/tminaorg/brzaguza/src/search"
)

func Search(c *gin.Context, config *config.Config, db cache.DB) {
	var query, pages, deepSearch string

	switch c.Request.Method {
	case "", "GET":
		{
			query = c.Query("q")
			pages = c.DefaultQuery("pages", "1")
			deepSearch = c.DefaultQuery("deep", "false")
		}
	case "POST":
		{
			query = c.PostForm("q")
			pages = c.DefaultPostForm("pages", "1")
			deepSearch = c.DefaultPostForm("deep", "false")
		}
	}

	maxPages, err := strconv.Atoi(pages)
	if err != nil {
		log.Error().Err(err).Msgf("cannot convert \"%v\" to int, reverting to default value of 1", pages)
		maxPages = 1
	}

	visitPages := false
	if deepSearch != "false" {
		log.Trace().Msgf("doing a deep search because deep is: %v", deepSearch)
		visitPages = true
	}

	options := engines.Options{
		MaxPages:   maxPages,
		VisitPages: visitPages,
	}

	var results []result.Result
	db.Get(query, &results)
	if results != nil {
		log.Debug().Msgf("Found results for query (%v) in cache", query)
	} else {
		log.Debug().Msg("Nothing found in cache, doing a clean search")
		results = search.PerformSearch(query, options, config)
		defer db.Set(query, results)
	}

	if resultsJson, err := json.Marshal(results); err != nil {
		log.Error().Err(err).Msg("failed marshalling results")
		c.String(http.StatusInternalServerError, "")
	} else {
		c.String(http.StatusOK, string(resultsJson))
	}
}
