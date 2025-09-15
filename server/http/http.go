package http

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v3"
)

type HTTP struct {
	App  *fiber.App
	Host string
}

func NewHTTP(host string) *HTTP {
	return &HTTP{App: fiber.New(fiber.Config{}), Host: host}
}

func (http *HTTP) Start() {
	log.Println("Starting HTTP handler")
	go func() {
		publicRouter(http.App)
		if err := http.App.Listen(http.Host); err != nil {
			log.Fatal(err)
		}
	}()
}

func (http *HTTP) Stop() {
	log.Println("Stopping HTTP handler")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := http.App.ShutdownWithContext(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("HTTP handler stopped")
}
