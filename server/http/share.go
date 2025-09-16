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
		file, err := fs.FindFileByShare(db, shid)
		if err != nil {
			log.Println(err)
			return c.SendStatus(fiber.StatusNotFound)
		}
		c.Set("Content-Type", file.Meta.ContentType)
		return c.Send(file.Data)
	})

	share.Post("/:fid", func(c fiber.Ctx) error {
		fid := c.Params("fid")
		shareID, err := fs.CreateShare(db, fid)
		if err != nil {
			log.Println(err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.SendString(shareID)
	})
}
