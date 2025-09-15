package http

import (
	"cloud-server/conf"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/redirect"
)

func publicRouter(api fiber.Router) {
	log.Println("Starting HTTP handler")
	pub := api.Group("/")

	pub.All("/*", redirect.New(redirect.Config{
		Rules: map[string]string{
			"/*": conf.GlobalConf.WebClient.Host + fmt.Sprintf(":%d", conf.GlobalConf.WebClient.Port),
		},
	}))
}
