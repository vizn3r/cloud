package http

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
)

func publicRouter(api fiber.Router) {
	pub := api.Group("/")

	pub.Get("/*", static.New("./public"))
}
