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

type Router struct {
	router *graceful.Graceful
	config *config.Config
}

func New(config *config.Config) (*Router, error) {
	router, err := graceful.Default(graceful.WithAddr(fmt.Sprintf(":%v", config.Server.Port)))
	return &Router{router: router, config: config}, err
}

func (r *Router) addCors() {
	r.router.Use(cors.New(cors.Config{
		AllowOrigins:     r.config.Server.FrontendUrls,
		AllowMethods:     []string{"HEAD", "GET", "POST"},
		AllowHeaders:     []string{"Origin", "X-Requested-With", "Content-Length", "Content-Type", "Accept"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))
}

func (r *Router) runWithContext(ctx context.Context, stopRouter context.CancelFunc) {
	if err := r.router.RunWithContext(ctx); err != context.Canceled {
		log.Error().Msgf("Failed starting router: %v", err)
	} else if err != nil {
		log.Info().Msgf("Stopping router...")
		stopRouter()
		r.router.Close()
		log.Debug().Msgf("Successfully stopped router")
	}
}

func (r *Router) Start(db cache.DB) {
	// signal interrupt
	ctx, stopRouter := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	// CORS
	r.addCors()

	// health
	r.router.GET("/healthz", HealthCheck)

	// search
	r.router.GET("/search", func(c *gin.Context) {
		Search(c, r.config, db)
	})
	r.router.POST("/search", func(c *gin.Context) {
		Search(c, r.config, db)
	})

	// startup
	r.runWithContext(ctx, stopRouter)
}
