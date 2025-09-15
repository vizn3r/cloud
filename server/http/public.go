package http

import (
	"cloud-server/conf"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/redirect"
)

func publicRouter(api fiber.Router) {
	log.Println("Starting HTTP handler")
	pub := api.Group("/")

	pub.Get("/*", redirect.New(redirect.Config{
		Rules: map[string]string{
			"/*": conf.GlobalConf.WebClientHost,
		},
	}))
}
