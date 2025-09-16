package auth

import (
	"cloud-server/db"
	"database/sql"
	"log"
	"time"

	"github.com/gofiber/fiber/v3"
)

// RequireFileOwnership middleware checks if the authenticated user owns the requested file
func RequireFileOwnership(data *db.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		log.Println("Checking file ownership")
		// Get file ID from route parameter
		fileID := c.Params("fid")
		if fileID == "" {
			return c.Status(fiber.StatusBadRequest).SendString("File ID required")
		}

		// Get user ID from authentication context
		userID, ok := c.Locals("userID").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).SendString("Authentication required")
		}

		// Query database to check file ownership
		var ownerID string
		var uploadedAt time.Time
		var updatedAt sql.NullTime // Use nullable type for updated_at which can be NULL
		err := data.Connection.QueryRow(db.Q_FILE_FIND_BY_ID, fileID).Scan(&ownerID, &uploadedAt, &updatedAt)
		if err != nil {
			log.Println("File ownership query failed:", err)
			return c.Status(fiber.StatusNotFound).SendString("File not found")
		}

		// Check if the authenticated user owns the file
		if ownerID != userID {
			return c.Status(fiber.StatusForbidden).SendString("Access denied - you don't own this file")
		}

		log.Println("Ownership OK")
		return c.Next()
	}
}
