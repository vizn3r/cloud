package http

import (
	"cloud-server/auth"
	"cloud-server/db"
	"cloud-server/fs"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
)

func shareRouter(api fiber.Router, db *db.DB) {
	share := api.Group("/share")

	share.Get("/:shid", func(c fiber.Ctx) error {
		shid := c.Params("shid")
		file, err := fs.FindFileByShare(db, shid)
		if err != nil {
			log.Error(err)
			return c.SendStatus(fiber.StatusNotFound)
		}
		c.Set("Content-Type", file.Meta.ContentType)
		return c.Send(file.Data)
	})

	share.Post("/:fid", auth.RequireToken(db), func(c fiber.Ctx) error {
		durationStr := c.Get("X-Share-Duration", "1440")
		duration, err := strconv.Atoi(durationStr)
		if err != nil || duration < 1 || duration > 10080 { // 1 min to 1 week max
			duration = 1440
		}
		fid := c.Params("fid")
		shareID, err := fs.CreateShare(db, fid, time.Duration(duration)*time.Minute)
		if err != nil {
			log.Error(err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.SendString(shareID)
	})
}
