package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hearchco/hearchco/src/cache"
	"github.com/hearchco/hearchco/src/config"
	"github.com/hearchco/hearchco/src/search/category"
	"github.com/hearchco/hearchco/src/search/engines"
)

func setupRoutes(app *fiber.App, db cache.DB, ttlConf config.TTL, settings map[engines.Name]config.Settings, categories map[category.Name]config.Category) {
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return HealthCheck(c)
	})

	app.Get("/search", func(c *fiber.Ctx) error {
		return Search(c, db, ttlConf, settings, categories)
	})

	app.Post("/search", func(c *fiber.Ctx) error {
		return Search(c, db, ttlConf, settings, categories)
	})
}
