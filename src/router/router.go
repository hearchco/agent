package router

import (
	"context"
	"fmt"
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
	ctx, stopRouter := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	// create router with configured port
	router, err := graceful.Default(graceful.WithAddr(fmt.Sprintf(":%v", config.Server.Port)))
	if err != nil {
		log.Error().Msgf("Failed creating a router: %v", err)
		return
	}

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
	if err := router.RunWithContext(ctx); err != nil {
		log.Error().Msgf("Failed starting router: %v", err)
	} else if err != context.Canceled {
		log.Info().Msgf("Stopping router...")
		stopRouter()
		router.Close()
		log.Info().Msgf("Successfully stopped router")
	}
}
