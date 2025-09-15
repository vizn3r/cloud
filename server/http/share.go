package http

import (
	"cloud-server/fs"
	"log"

	"github.com/gofiber/fiber/v3"
)

func handleShare(api fiber.Router) {
	share := api.Group("/share")

	share.Get("/:shid", func(c fiber.Ctx) error {
		shid := c.Params("shid")
		file, err := fs.FindShare(shid)
		if err != nil {
			log.Println(err)
			c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.Send(file.Data)
	})
}
