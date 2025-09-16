package http

import (
	"cloud-server/auth"
	"cloud-server/db"
	"cloud-server/user"
	"time"

	"github.com/gofiber/fiber/v3"
)

func userRouter(api fiber.Router, db *db.DB) {
	usr := api.Group("/user")

	usr.Post("/register", func(c fiber.Ctx) error {
		type RegisterRequest struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		var req RegisterRequest
		if err := c.Bind().Body(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request",
			})
		}

		userID, err := user.CreateUser(db, req.Email, req.Password)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Create session and return token (auto-login after registration)
		token, err := user.CreateSession(db, userID, 24*time.Hour)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create session",
			})
		}

		return c.JSON(fiber.Map{
			"message": "User created successfully",
			"user_id": userID,
			"token": token,
		})
	})

	usr.Post("/login", func(c fiber.Ctx) error {
		type LoginRequest struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		var req LoginRequest
		if err := c.Bind().Body(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request",
			})
		}

		userID, err := user.AuthenticateUser(db, req.Email, req.Password)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid credentials",
			})
		}

		token, err := user.CreateSession(db, userID, 24*time.Hour) // 24 hour session
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create session",
			})
		}

		return c.JSON(fiber.Map{
			"message": "Login successful",
			"token":   token,
			"user_id": userID,
		})
	})

	usr.Get("/me", auth.RequireToken(db), func(c fiber.Ctx) error {
		userID := c.Locals("userID").(string)

		userInfo, err := user.GetUserByID(db, userID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}

		return c.JSON(fiber.Map{
			"id":         userInfo.ID,
			"email":      userInfo.Email,
			"created_at": userInfo.CreatedAt,
		})
	})

	usr.Get("/files", auth.RequireToken(db), func(c fiber.Ctx) error {
		userID := c.Locals("userID").(string)

		files, err := user.GetUserFiles(db, userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to get files",
			})
		}

		return c.JSON(fiber.Map{
			"files": files,
		})
	})
}
