package main

import (
	"fmt"
	"log"
	"webclient/conf"

	_ "embed"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
)

//go:embed webclient.json
var config string

func main() {
	if err := conf.LoadFromBytes([]byte(config)); err != nil {
		log.Fatal(err)
	}

	app := fiber.New(fiber.Config{})

	app.Use("/*", static.New("./public"))

	app.Listen(fmt.Sprintf(":%d", conf.GlobalConf.Port), fiber.ListenConfig{})
}
