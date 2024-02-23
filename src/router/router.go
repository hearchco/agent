package router

import (
	"context"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/graceful"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/hearchco/logger"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/hearchco/hearchco/src/cache"
	"github.com/hearchco/hearchco/src/config"
)

type RouterWrapper struct {
	rtr  *graceful.Graceful
	conf config.Config
}

func New(conf config.Config, verbosity int8, lgr zerolog.Logger) (RouterWrapper, error) {
	// set verbosity to release mode if log level is INFO
	if verbosity == 0 {
		gin.SetMode(gin.ReleaseMode)
	}

	// create new gin engine with recovery middleware and zerolog logger
	gengine := gin.New()
	gengine.Use(gin.Recovery())
	gengine.Use(logger.SetLogger(logger.WithLogger(func(c *gin.Context, l zerolog.Logger) zerolog.Logger {
		return lgr.With().
			Int("status", c.Writer.Status()).
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Str("ip", c.ClientIP()).
			Logger()
	}), logger.WithDefaultFieldsDisabled(), logger.WithLatency(), logger.WithSkipPath([]string{"/health", "/healthz"})))

	// add CORS middleware
	log.Debug().
		Str("url", conf.Server.FrontendUrl).
		Msg("Using CORS")
	gengine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{conf.Server.FrontendUrl},
		AllowMethods:     []string{"HEAD", "GET", "POST"},
		AllowHeaders:     []string{"Origin", "X-Requested-With", "Content-Length", "Content-Type", "Accept"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	// create new graceful engine with config port
	rtr, err := graceful.New(gengine, graceful.WithAddr(":"+strconv.Itoa(conf.Server.Port)))
	if err != nil {
		log.Error().
			Err(err).
			Msg("router.New(): failed creating new graceful router")
	}

	return RouterWrapper{rtr: rtr, conf: conf}, err
}

func (rw RouterWrapper) runWithContext(ctx context.Context) error {
	log.Info().
		Int("port", rw.conf.Server.Port).
		Msg("Started router")
	if err := rw.rtr.RunWithContext(ctx); err != context.Canceled {
		log.Error().
			Err(err).
			Msg("router.runWithContext(): failed starting router")
		return err
	} else if err != nil {
		log.Info().Msg("Stopping router...")
		rw.rtr.Close()
		log.Debug().Msg("Successfully stopped router")
	}
	return nil
}

func (rw RouterWrapper) Start(ctx context.Context, db cache.DB, serveProfiler bool) error {
	// health(z)
	rw.rtr.GET("/health", HealthCheck)
	rw.rtr.GET("/healthz", HealthCheck)

	// search
	rw.rtr.GET("/search", func(c *gin.Context) {
		err := Search(c, rw.conf, db)
		if err != nil {
			log.Error().Err(err).Msg("router.Start() (.GET): failed search")
		}
	})
	rw.rtr.POST("/search", func(c *gin.Context) {
		err := Search(c, rw.conf, db)
		if err != nil {
			log.Error().Err(err).Msg("router.Start() (.POST): failed search")
		}
	})

	if serveProfiler {
		pprof.Register(rw.rtr.Engine)
	}
	return rw.runWithContext(ctx)
}
