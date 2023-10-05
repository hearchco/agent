package router

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/graceful"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/tminaorg/brzaguza/src/cache"
	"github.com/tminaorg/brzaguza/src/config"
)

func Setup(config *config.Config, db cache.DB) {
	// signal interrupt
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// create router
	router, err := graceful.Default()
	if err != nil {
		log.Error().Msgf("Failed creating a router: %v", err)
		return
	}
	defer router.Close()

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
	//router.Run(fmt.Sprintf(":%v", config.Server.Port))
	if err := router.RunWithContext(ctx); err != nil {
		log.Error().Msgf("Failed creating a router: %v", err)
	} else if err != context.Canceled {
		log.Info().Msgf("Stopped router")
	}
}
