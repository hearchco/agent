package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/tminaorg/brzaguza/src/config"
	"github.com/tminaorg/brzaguza/src/search"
)

func SetupRouter(config *config.Config) {
	router := gin.Default()

	router.GET("/search", func(c *gin.Context) {
		query := c.Query("q")

		pages := c.Query("pages")
		maxPages := 1
		if pages != "" {
			tmpMapPages, err := strconv.Atoi(pages)
			if err != nil {
				log.Error().Err(err).Msgf("cannot convert maxPages=%v to int, reverting to default value of 1", pages)
				maxPages = tmpMapPages
			}
		}

		deepSearch := c.Query("deep")
		visit := false
		if deepSearch != "" {
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

	router.Run(fmt.Sprintf(":%v", strconv.Itoa(config.Server.Port)))
}