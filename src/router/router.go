package router

import (
	"context"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/graceful"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/hearchco/hearchco/src/cache"
	"github.com/hearchco/hearchco/src/config"
)

type RouterWrapper struct {
	router *graceful.Graceful
	config *config.Config
}

func New(config *config.Config, verbosity int8) (*RouterWrapper, error) {
	if verbosity == 0 {
		gin.SetMode(gin.ReleaseMode)
	}
	router, err := graceful.Default(graceful.WithAddr(":" + strconv.Itoa(config.Server.Port)))
	return &RouterWrapper{router: router, config: config}, err
}

func (rw *RouterWrapper) addCors() {
	log.Debug().
		Str("URL", rw.config.Server.FrontendUrl).
		Msg("Using CORS")
	rw.router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{rw.config.Server.FrontendUrl},
		AllowMethods:     []string{"HEAD", "GET", "POST"},
		AllowHeaders:     []string{"Origin", "X-Requested-With", "Content-Length", "Content-Type", "Accept"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))
}

func (rw *RouterWrapper) runWithContext(ctx context.Context) {
	if err := rw.router.RunWithContext(ctx); err != context.Canceled {
		log.Error().Err(err).Msg("router.runWithContext(): failed starting router")
	} else if err != nil {
		log.Info().Msg("Stopping router...")
		rw.router.Close()
		log.Debug().Msg("Successfully stopped router")
	}
}

func (rw *RouterWrapper) Start(ctx context.Context, db cache.DB, serveProfiler bool) {
	// CORS
	rw.addCors()

	// health
	rw.router.GET("/healthz", HealthCheck)

	// search
	rw.router.GET("/search", func(c *gin.Context) {
		err := Search(c, rw.config, db)
		if err != nil {
			log.Error().Err(err).Msgf("router.Start() (.GET): failed search")
		}
	})
	rw.router.POST("/search", func(c *gin.Context) {
		err := Search(c, rw.config, db)
		if err != nil {
			log.Error().Err(err).Msgf("router.Start() (.POST): failed search")
		}
	})

	if serveProfiler {
		pprof.Register(rw.router.Engine)
	}
	rw.runWithContext(ctx)
}
