package router

import (
	"github.com/gofiber/fiber/v2"

	"github.com/hearchco/hearchco/src/cache"
	"github.com/hearchco/hearchco/src/config"
)

func setupRoutes(app *fiber.App, db cache.DB, conf config.Config) {
	app.Get("/search", func(c *fiber.Ctx) error {
		return Search(c, db, conf.Server.Cache.TTL, conf.Settings, conf.Categories, conf.Server.Proxy.Salt)
	})

	app.Post("/search", func(c *fiber.Ctx) error {
		return Search(c, db, conf.Server.Cache.TTL, conf.Settings, conf.Categories, conf.Server.Proxy.Salt)
	})

	app.Get("/proxy", func(c *fiber.Ctx) error {
		return Proxy(c, conf.Server.Proxy.Salt, conf.Server.Proxy.Timeout)
	})
}
