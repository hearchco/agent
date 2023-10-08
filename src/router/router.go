package router

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/graceful"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/tminaorg/brzaguza/src/cache"
	"github.com/tminaorg/brzaguza/src/config"
)

type RouterWrapper struct {
	router *graceful.Graceful
	config *config.Config
}

func New(config *config.Config, verbosity int8) (*RouterWrapper, error) {
	if verbosity == 0 {
		gin.SetMode(gin.ReleaseMode)
	}
	router, err := graceful.Default(graceful.WithAddr(fmt.Sprintf(":%v", config.Server.Port)))
	return &RouterWrapper{router: router, config: config}, err
}

func (rw *RouterWrapper) addCors() {
	rw.router.Use(cors.New(cors.Config{
		AllowOrigins:     rw.config.Server.FrontendUrls,
		AllowMethods:     []string{"HEAD", "GET", "POST"},
		AllowHeaders:     []string{"Origin", "X-Requested-With", "Content-Length", "Content-Type", "Accept"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))
}

func (rw *RouterWrapper) runWithContext(ctx context.Context) {
	if err := rw.router.RunWithContext(ctx); err != context.Canceled {
		log.Error().Msgf("Failed starting router: %v", err)
	} else if err != nil {
		log.Info().Msgf("Stopping router...")
		rw.router.Close()
		log.Debug().Msgf("Successfully stopped router")
	}
}

func (rw *RouterWrapper) Start(ctx context.Context, db cache.DB, serveProfiler bool) {
	// CORS
	rw.addCors()

	// health
	rw.router.GET("/healthz", HealthCheck)

	// search
	rw.router.GET("/search", func(c *gin.Context) {
		Search(c, rw.config, db)
	})
	rw.router.POST("/search", func(c *gin.Context) {
		Search(c, rw.config, db)
	})

	if serveProfiler {
		pprof.Register(rw.router.Engine)
	}
	rw.runWithContext(ctx)
}
