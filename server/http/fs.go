package http

import (
	"cloud-server/fs"
	"fmt"
	"log"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gofiber/fiber/v3"
)

func fsRouter(api fiber.Router) {
	files := api.Group("/file")

	files.Get("/:fid", func(c fiber.Ctx) error {
		fid := c.Params("fid")
		log.Println("Requesting file: ", fid)
		data := fs.FindFile(fid)
		mimeType := mimetype.Detect(data)
		c.Set("Content-Type", mimeType.String())
		return c.Send(data)
	})

	files.Post("/", func(c fiber.Ctx) error {
		file, err := c.FormFile("file")
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		opened, err := file.Open()
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		defer opened.Close()

		data := make([]byte, file.Size)
		_, err = opened.Read(data)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		id := fs.SaveFile(data)
		log.Println("Uploaded file: ", id)
		return c.SendString(fmt.Sprintf("%d", id))
	})
}
