package router

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tminaorg/brzaguza/src/cache"
	"github.com/tminaorg/brzaguza/src/config"
)

var router = gin.Default()

func Setup(config *config.Config, db cache.DB) {
	// CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     config.Server.FrontendUrls,
		AllowMethods:     []string{"HEAD", "GET", "POST"},
		AllowHeaders:     []string{"Origin", "X-Requested-With", "Content-Length", "Content-Type", "Accept"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	// health
	router.GET("/healthz", HealthCheck)

	// search
	router.GET("/search", func(c *gin.Context) {
		Search(c, config, db)
	})
	router.POST("/search", func(c *gin.Context) {
		Search(c, config, db)
	})

	// startup
	router.Run(fmt.Sprintf(":%v", config.Server.Port))
}
