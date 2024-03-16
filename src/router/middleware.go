package router

import (
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func setupMiddlewares(app *fiber.App, lgr zerolog.Logger, frontendUrls []string) {
	// use recovery middleware
	app.Use(recover.New())

	// use zerolog middleware
	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger:   &lgr,
		SkipURIs: []string{"/livez", "/readyz"},
		Fields: []string{
			fiberzerolog.FieldIP,
			fiberzerolog.FieldLatency,
			fiberzerolog.FieldStatus,
			fiberzerolog.FieldMethod,
			fiberzerolog.FieldPath,
			fiberzerolog.FieldError,
		},
	}))

	// use gzip, deflate and brotli middleware
	app.Use(compress.New())

	// use favicon ignore middleware
	app.Use(favicon.New())

	// use CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "",
		AllowOriginsFunc: CheckOrigin(frontendUrls),
	}))

	log.Debug().
		Strs("url", frontendUrls).
		Msg("Using CORS")

	// use liveness and readiness middleware
	app.Use(healthcheck.New())
}
