package http

import (
	"cloud-server/db"
	"cloud-server/fs"
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gofiber/fiber/v3"
)

func fsRouter(api fiber.Router, db *db.DB) {
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
			return c.SendStatus(fiber.StatusNotFound)
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
		// File upload size limit: 50MB
		const maxUploadSize = 50 * 1024 * 1024

		var data []byte
		ogName := ""
		file, err := c.FormFile("file")
		if err == nil {
			// Validate file size
			if file.Size > maxUploadSize {
				return c.Status(fiber.StatusRequestEntityTooLarge).SendString("File too large")
			}

			// Validate file name
			if !isSafeFilename(file.Filename) {
				return c.Status(fiber.StatusBadRequest).SendString("Invalid filename")
			}

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
			// Validate raw body size
			if len(c.BodyRaw()) > maxUploadSize {
				return c.Status(fiber.StatusRequestEntityTooLarge).SendString("Request too large")
			}
			data = c.BodyRaw()
		}

		mimeType := mimetype.Detect(data)

		newFile := fs.File{
			Meta: fs.FileMeta{
				UploadName:  c.Get("X-Original-Filename", ogName),
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

func isSafeFilename(filename string) bool {
	// Prevent path traversal and malicious filenames
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
		return false
	}

	// Prevent potentially dangerous extensions
	blacklist := []string{".exe", ".bat", ".cmd", ".sh", ".php", ".py", ".js", ".html"}
	for _, ext := range blacklist {
		if strings.HasSuffix(strings.ToLower(filename), ext) {
			return false
		}
	}

	return true
}
