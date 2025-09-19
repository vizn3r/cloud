package auth

import (
	"database/sql"
	"time"

	"cloud-server/db"

	"github.com/gofiber/fiber/v3"
)

// RequireFileOwnership middleware checks if the authenticated user owns the requested file
func RequireFileOwnership(data *db.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		log.Print("Owner check...")
		// Get file ID from route parameter
		fileID := c.Params("fid")
		if fileID == "" {
			log.Print("File ID FAIL")
			return c.Status(fiber.StatusBadRequest).SendString("File ID required")
		}

		// Get user ID from authentication context
		userID, ok := c.Locals("userID").(string)
		if !ok {
			log.Print("Auth FAIL")
			return c.Status(fiber.StatusUnauthorized).SendString("Authentication required")
		}

		// Query database to check file ownership
		var ownerID string
		var uploadedAt time.Time
		var updatedAt sql.NullTime
		err := data.Connection.QueryRow(db.Q_FILE_FIND_BY_ID, fileID).Scan(&ownerID, &uploadedAt, &updatedAt)
		if err != nil {
			log.Error("File ownership query failed: ", err)
			return c.Status(fiber.StatusNotFound).SendString("File not found")
		}

		// Check if the authenticated user owns the file
		if ownerID != userID {
			log.Print("Owner ID FAIL")
			return c.Status(fiber.StatusForbidden).SendString("Access denied - you don't own this file")
		}

		log.Print("Ownership OK")
		return c.Next()
	}
}
