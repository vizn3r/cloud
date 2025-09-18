package http

import (
	"cloud-server/auth"
	"cloud-server/db"
	"cloud-server/fs"
	"encoding/json"
	"fmt"
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
		log.Print("Requesting file: ", fid)

		file, err := fs.FindFile(fid)
		if os.IsNotExist(err) {
			log.Error(err)
			return c.SendStatus(fiber.StatusNotFound)
		} else if err != nil {
			log.Error(err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		c.Set("Content-Type", file.Meta.ContentType)
		return c.Send(file.Data)
	})

	files.Get("/:fid/data", auth.RequireToken(data), auth.RequireFileOwnership(data), func(c fiber.Ctx) error {
		fid := c.Params("fid")
		log.Print("Requesting file: ", fid)

		file, err := fs.FindFile(fid)
		if os.IsNotExist(err) {
			log.Error(err)
			return c.SendStatus(fiber.StatusNotFound)
		} else if err != nil {
			log.Error(err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		c.Set("Content-Type", "text/json")

		data, err := json.Marshal(file.Meta)
		if err != nil {
			log.Error(err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.Send(data)
	})

	files.Get("/:fid/thumbnail", auth.RequireToken(data), auth.RequireFileOwnership(data), func(c fiber.Ctx) error {
		fid := c.Params("fid")
		log.Print("Requesting thumbnail for file: ", fid)

		file, err := fs.FindFile(fid)
		if os.IsNotExist(err) {
			log.Error(err)
			return c.SendStatus(fiber.StatusNotFound)
		} else if err != nil {
			log.Error(err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		// Check if it's an image
		if strings.HasPrefix(file.Meta.ContentType, "image/") {
			// Return the actual image as thumbnail
			c.Set("Content-Type", file.Meta.ContentType)
			return c.Send(file.Data)
		}

		// For non-image files, generate an icon based on file extension
		c.Set("Content-Type", "image/svg+xml")

		// Extract file extension from upload name
		extension := "file"
		if file.Meta.UploadName != "" {
			parts := strings.Split(file.Meta.UploadName, ".")
			if len(parts) > 1 {
				extension = parts[len(parts)-1]
			}
		}

		// Generate SVG icon with file extension
		iconSVG := fileIconSVG(extension)
		return c.SendString(iconSVG)
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
			if !checkFileName(file.Filename) {
				return c.Status(fiber.StatusBadRequest).SendString("Invalid filename")
			}

			ogName = file.Filename
			opened, err := file.Open()
			if err != nil {
				log.Error(err)
				return c.SendStatus(fiber.StatusInternalServerError)
			}
			defer opened.Close()

			fileData = make([]byte, file.Size)
			_, err = opened.Read(fileData)
			if err != nil {
				log.Error(err)
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
			log.Error(err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		log.Print("Uploaded file:", id, "by user:", userID)
		return c.SendString(id)
	})

	files.Delete("/:fid", auth.RequireToken(data), auth.RequireFileOwnership(data), func(c fiber.Ctx) error {
		fid := c.Params("fid")

		err := fs.DeleteFile(fid, data.Connection)
		if err != nil {
			log.Error("Failed to delete file", fid, err)
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to delete file")
		}

		log.Print("Deleted file:", fid, "by user:", c.Locals("userID").(string))
		return c.SendStatus(fiber.StatusOK)
	})
}

func checkFileName(filename string) bool {
	return !strings.Contains(filename, "..") && !strings.Contains(filename, "/") && !strings.Contains(filename, "\\")
}

func fileIconSVG(extension string) string {
	// Simple SVG icon with file extension text
	return fmt.Sprintf(`<svg width="64" height="64" xmlns="http://www.w3.org/2000/svg">
		<rect width="64" height="64" fill="#3b82f6" rx="8"/>
		<text x="32" y="32" font-family="Arial, sans-serif" font-size="12" fill="white"
			text-anchor="middle" dominant-baseline="middle">%s</text>
	</svg>`, strings.ToUpper(extension))
}
