package auth

import (
	"os"
	"strings"

	"cloud-server/db"
	"cloud-server/logger"
	"cloud-server/user"

	"github.com/gofiber/fiber/v3"
)

func IsTest() bool {
	return os.Getenv("TEST") == "true"
}

var log = logger.New("AUTH", logger.Red)

func RequireToken(db *db.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		if IsTest() {
			log.Warn("Test mode - auth disabled")
			c.Locals("userID", "test-user-uuid")
			return c.Next()
		}
		log.Print("Token check...")
		authHeader := c.Get("Authorization")

		if !strings.HasPrefix(authHeader, "Bearer ") {
			log.Print("Auth FAIL")
			return c.Status(fiber.StatusUnauthorized).SendString("Authorization required")
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			log.Print("Token FAIL")
			return c.Status(fiber.StatusUnauthorized).SendString("Invalid token")
		}

		userID, err := user.ValidateSession(db, token)
		if err != nil {
			log.Print("Session FAIL")
			return c.Status(fiber.StatusUnauthorized).SendString("Invalid session")
		}

		c.Locals("userID", userID)
		log.Print("Access OK")
		return c.Next()
	}
}
