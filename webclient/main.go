package main

import (
	"fmt"
	"log"
	"webclient/conf"

	_ "embed"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/static"
)

//go:embed webclient.json
var config string

func main() {
	if err := conf.LoadFromBytes([]byte(config)); err != nil {
		log.Fatal(err)
	}

	app := fiber.New(fiber.Config{})

	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${ip} | ${status} | ${method} | ${path} | ${latency}\n",
		TimeFormat: "02-Jan-2006 15:04:05",
		TimeZone:   "Local",
	}))

	app.Use("/*", static.New("./public"))

	if err := app.Listen(fmt.Sprintf(":%d", conf.GlobalConf.Port), fiber.ListenConfig{}); err != nil {
		log.Fatal(err)
	}
}
