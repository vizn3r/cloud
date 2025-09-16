package main

import (
	"fmt"
	"net/http"
	"strings"
	"webclient/conf"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
)

func main() {
	conf.LoadConfig("webclient.json")

	app := fiber.New(fiber.Config{})

	// Authentication middleware for /app routes
	app.Use("/app/*", func(c fiber.Ctx) error {
		// Check if user is authenticated by verifying session token
		authHeader := c.Get("Authorization")

		if authHeader == "" {
			// Check for token in query parameter (for initial redirect)
			token := c.Query("token")
			if token == "" {
				// No token found, redirect to login
				return c.Redirect().To("/login")
			}
			// Set Authorization header from query parameter
			c.Request().Header.Set("Authorization", "Bearer "+token)
			authHeader = "Bearer " + token
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Redirect().To("/login")
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			return c.Redirect().To("/login")
		}

		// Verify token with auth server
		client := &http.Client{}
		req, _ := http.NewRequest("GET", "http://localhost:8080/user/me", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := client.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			return c.Redirect().To("/login")
		}
		defer resp.Body.Close()

		// Token is valid, continue to serve the app
		return c.Next()
	})

	app.Use("/*", static.New("./public"))

	app.Listen(fmt.Sprintf(":%d", conf.GlobalConf.Port), fiber.ListenConfig{})
}
