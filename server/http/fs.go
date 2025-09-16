package http

import (
	"cloud-server/auth"
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

func fsRouter(api fiber.Router, data *db.DB) {
	files := api.Group("/file")

	files.Get("/:fid", auth.RequireToken(data), auth.RequireFileOwnership(data), func(c fiber.Ctx) error {
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

	files.Get("/:fid/data", auth.RequireToken(data), auth.RequireFileOwnership(data), func(c fiber.Ctx) error {
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

	files.Post("/", auth.RequireToken(data), func(c fiber.Ctx) error {
		// File upload size limit: 50MB
		const maxUploadSize = 50 * 1024 * 1024

		var fileData []byte
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

			fileData = make([]byte, file.Size)
			_, err = opened.Read(fileData)
			if err != nil {
				log.Println(err)
				return c.SendStatus(fiber.StatusInternalServerError)
			}
		} else {
			// Validate raw body size
			if len(c.BodyRaw()) > maxUploadSize {
				return c.Status(fiber.StatusRequestEntityTooLarge).SendString("Request too large")
			}
			fileData = c.BodyRaw()
		}

		mimeType := mimetype.Detect(fileData)

		newFile := fs.File{
			Meta: fs.FileMeta{
				UploadName:  c.Get("X-Original-Filename", ogName),
				Size:        uint64(len(fileData)),
				UploadedAt:  time.Now(),
				ContentType: mimeType.String(),
			},
			Data: fileData,
		}

		// Get user ID for file ownership
		userID := c.Locals("userID").(string)
		id, err := newFile.SaveFile(data.Connection, userID)
		if err != nil {
			log.Println(err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		log.Println("Uploaded file:", id, "by user:", userID)
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
