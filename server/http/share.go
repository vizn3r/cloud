package http

import (
	"cloud-server/db"
	"cloud-server/fs"
	"log"

	"github.com/gofiber/fiber/v3"
)

func shareRouter(api fiber.Router, db *db.DB) {
	share := api.Group("/share")

	share.Get("/:shid", func(c fiber.Ctx) error {
		shid := c.Params("shid")
		share := fs.NewShare(shid, db)
		file, err := share.FindShare()
		if err != nil {
			log.Println(err)
			c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.Send(file.Data)
	})

	share.Post("/:fid", func(c fiber.Ctx) error {
		fid := c.Params("fid")
		share := &fs.Share{
			FileUUID: fid,
		}
		newShare, err := share.SaveShare()
		if err != nil {
			log.Println(err)
			c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.SendString(newShare.UUID)
	})
}
