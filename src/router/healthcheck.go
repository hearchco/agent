package router

import "github.com/gofiber/fiber/v2"

func HealthCheck(c *fiber.Ctx) error {
	// TODO: check db health
	// if err -> return c.Send(error)
	return c.SendString("OK")
}
