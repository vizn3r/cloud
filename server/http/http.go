package http

import (
	"cloud-server/conf"
	"cloud-server/db"
	"cloud-server/logger"
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

var log = logger.New("HTTP", logger.Green)

type HTTP struct {
	App  *fiber.App
	Host string
	db   *db.DB
	done chan struct{}
}

func NewHTTP(host string, db *db.DB) *HTTP {
	return &HTTP{App: fiber.New(fiber.Config{}), Host: host, db: db}
}

func (http *HTTP) Start() {
	http.done = make(chan struct{})
	log.Info("Starting HTTP handler")
	go func() {
		http.App.Use(func(c fiber.Ctx) error {
			log.Print("Request: ", c.Method(), c.Path())
			return c.Next()
		})

		http.App.Use(cors.New(cors.Config{
			AllowOrigins: []string{
				fmt.Sprintf("%s:%d", conf.GlobalConf.WebClient.Host, conf.GlobalConf.WebClient.Port),
				"https://cloud.vizn3r.eu",
				"http://cloud.vizn3r.eu",
				"https://cloudapi.vizn3r.eu",
				"http://cloudapi.vizn3r.eu",
			},
			AllowMethods:     []string{"GET,POST,PUT,DELETE,OPTIONS"},
			AllowHeaders:     []string{"Content-Type", "Authorization", "X-Original-Filename", "X-Share-Duration"},
			AllowCredentials: true,
			MaxAge:           300,
		}))
		v1 := http.App.Group("/v1")
		fsRouter(v1, http.db)
		shareRouter(v1, http.db)
		userRouter(v1, http.db)
		publicRouter(v1)
		if err := http.App.Listen(http.Host); err != nil {
			log.Fatal(err)
		}
		close(http.done)
	}()
}

func (http *HTTP) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := http.App.ShutdownWithContext(ctx); err != nil {
		log.Fatal(err)
	}
	log.Warn("HTTP handler stopped")
	<-http.done
	log.Close()
}
