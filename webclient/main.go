package main

import (
	"fmt"
	"webclient/conf"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
)

func main() {
	conf.LoadConfig("webclient.json")

	app := fiber.New(fiber.Config{})

	app.Use("/*", static.New("./public"))

	app.Listen(fmt.Sprintf(":%d", conf.GlobalConf.Port), fiber.ListenConfig{})
}
