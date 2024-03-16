package router

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/hearchco/hearchco/src/anonymize"
)

func Proxy(c *fiber.Ctx, salt string, timeout time.Duration) error {
	url := c.Query("url")
	hash := c.Query("hash")

	if url == "" || hash == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Message: "\"url\" and \"hash\" are required",
		})
	}

	// check if hash is valid
	if !anonymize.CheckHash(hash, url, salt) {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Message: "Invalid hash",
			Value:   hash,
		})
	}

	// wait for maximum of timeout
	if err := proxy.DoTimeout(c, url, timeout); err != nil {
		return err
	}

	// remove server header from response
	c.Response().Header.Del(fiber.HeaderServer)

	return nil
}
