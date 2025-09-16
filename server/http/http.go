package http

import (
	"cloud-server/conf"
	"cloud-server/db"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

type HTTP struct {
	App  *fiber.App
	Host string
	db   *db.DB
}

func NewHTTP(host string, db *db.DB) *HTTP {
	return &HTTP{App: fiber.New(fiber.Config{}), Host: host, db: db}
}

func (http *HTTP) Start() {
	log.Println("Starting HTTP handler")
	go func() {
		http.App.Use(cors.New(cors.Config{
			AllowOrigins: []string{conf.GlobalConf.WebClient.Host + fmt.Sprintf(":%d", conf.GlobalConf.WebClient.Port)},
		}))
		http.App.Use(func(c fiber.Ctx) error {
			log.Println("New req: ", c.Request().String())
			return c.Next()
		})
		v1 := http.App.Group("/")
		fsRouter(v1, http.db)
		shareRouter(v1, http.db)
		publicRouter(v1)
		if err := http.App.Listen(http.Host); err != nil {
			log.Fatal(err)
		}
	}()
}

func (http *HTTP) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := http.App.ShutdownWithContext(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("HTTP handler stopped")
}
