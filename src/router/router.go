package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/tminaorg/brzaguza/src/config"
)

var router = gin.Default()

func Setup(config *config.Config) {
	// health
	router.GET("/healthz", HealthCheck)

	// search
	router.GET("/search", func(c *gin.Context) {
		Search(c, config)
	})
	router.POST("/search", func(c *gin.Context) {
		Search(c, config)
	})

	// startup
	router.Run(fmt.Sprintf(":%v", config.Server.Port))
}
