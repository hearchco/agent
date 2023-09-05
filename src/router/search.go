package router

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/tminaorg/brzaguza/src/config"
	"github.com/tminaorg/brzaguza/src/search"
)

func Search(config *config.Config) {
	searchRoute := router.Group("/search")

	searchRoute.GET("", func(c *gin.Context) {
		query := c.Query("q")

		pages := c.DefaultQuery("pages", "1")
		maxPages, err := strconv.Atoi(pages)
		if err != nil {
			log.Error().Err(err).Msgf("cannot convert \"%v\" to int, reverting to default value of 1", pages)
			maxPages = 1
		}

		deepSearch := c.DefaultQuery("deep", "false")
		visit := false
		if deepSearch != "false" {
			log.Trace().Msgf("doing a deep search because deep is: %v", deepSearch)
			visit = true
		}

		results := search.PerformSearch(query, maxPages, visit, config)

		if resultsJson, err := json.Marshal(results); err != nil {
			log.Error().Err(err).Msg("failed marshalling results")
			c.String(http.StatusInternalServerError, "")
		} else {
			c.String(http.StatusOK, string(resultsJson))
		}
	})

	searchRoute.POST("", func(c *gin.Context) {
		query := c.PostForm("q")

		pages := c.DefaultPostForm("pages", "1")
		maxPages, err := strconv.Atoi(pages)
		if err != nil {
			log.Error().Err(err).Msgf("cannot convert \"%v\" to int, reverting to default value of 1", pages)
			maxPages = 1
		}

		deepSearch := c.DefaultPostForm("deep", "false")
		visit := false
		if deepSearch != "false" {
			log.Trace().Msgf("doing a deep search because deep is: %v", deepSearch)
			visit = true
		}

		results := search.PerformSearch(query, maxPages, visit, config)

		if resultsJson, err := json.Marshal(results); err != nil {
			log.Error().Err(err).Msg("failed marshalling results")
			c.String(http.StatusInternalServerError, "")
		} else {
			c.String(http.StatusOK, string(resultsJson))
		}
	})
}
