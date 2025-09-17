package http

import (
	"github.com/gofiber/fiber/v3"
)

func publicRouter(api fiber.Router) {
	// Health check endpoint
	api.Get("/health", func(c fiber.Ctx) error {
		return c.SendString("OK")
	})
}
