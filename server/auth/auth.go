package auth

import (
	"cloud-server/db"
	"cloud-server/user"
	"log"
	"strings"

	"github.com/gofiber/fiber/v3"
)

func RequireToken(db *db.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		log.Println("Checking access token")
		authHeader := c.Get("Authorization")

		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).SendString("Authorization required")
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).SendString("Invalid token")
		}

		userID, err := user.ValidateSession(db, token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString("Invalid session")
		}

		c.Locals("userID", userID)
		log.Println("Token OK")
		return c.Next()
	}
}
