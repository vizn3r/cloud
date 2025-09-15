package http

import (
	"cloud-server/fs"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gofiber/fiber/v3"
)

func fsRouter(api fiber.Router) {
	files := api.Group("/file")

	files.Get("/:fid", func(c fiber.Ctx) error {
		fid := c.Params("fid")
		log.Println("Requesting file: ", fid)
		file, err := fs.FindFile(fid)
		if os.IsNotExist(err) {
			log.Println(err)
			return c.SendStatus(fiber.StatusNotFound)
		} else if err != nil {
			log.Println(err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		c.Set("Content-Type", file.Meta.ContentType)
		return c.Send(file.Data)
	})

	files.Get("/:fid/data", func(c fiber.Ctx) error {
		fid := c.Params("fid")
		log.Println("Requesting file: ", fid)
		file, err := fs.FindFile(fid)
		if os.IsNotExist(err) {
			log.Println(err)
			return c.SendStatus(fiber.StatusInternalServerError)
		} else if err != nil {
			log.Println(err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		c.Set("Content-Type", "text/json")

		data, err := json.Marshal(file.Meta)
		if err != nil {
			log.Println(err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.Send(data)
	})

	files.Post("/", func(c fiber.Ctx) error {
		var data []byte
		ogName := ""
		file, err := c.FormFile("file")
		if err == nil {
			ogName = file.Filename
			opened, err := file.Open()
			if err != nil {
				log.Println(err)
				return c.SendStatus(fiber.StatusInternalServerError)
			}
			defer opened.Close()

			data = make([]byte, file.Size)
			_, err = opened.Read(data)
			if err != nil {
				log.Println(err)
				return c.SendStatus(fiber.StatusInternalServerError)
			}
		} else {
			data = c.BodyRaw()
		}

		mimeType := mimetype.Detect(data)

		newFile := fs.File{
			Meta: fs.FileMeta{
				UploadName:  c.Get("X-Original-filename", ogName),
				Size:        uint64(len(data)),
				UploadedAt:  time.Now(),
				ContentType: mimeType.String(),
			},
			Data: data,
		}

		id, err := newFile.SaveFile()
		if err != nil {
			log.Println(err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		log.Println("Uploaded file: ", id)
		return c.SendString(id)
	})
}
