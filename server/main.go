package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"cloud-server/conf"
	"cloud-server/db"
	"cloud-server/http"
	"cloud-server/logger"
)

//go:embed server.json
var config string

//go:embed art.txt
var art string
var log = logger.New("MAIN", logger.Cyan)

func main() {
	log.Info("Loading config...")
	if err := conf.LoadFromBytes([]byte(config)); err != nil {
		log.Fatal(err)
	}

	if http.IsTest() {
		log.Warn("Running in test mode")
	}

	// Ensure temp directory exists
	if _, err := os.Stat("storage/temp"); os.IsNotExist(err) {
		err := os.MkdirAll("storage/temp", 0o700)
		if err != nil {
			log.Fatal(err)
		}
	}

	dbHost := db.NewDB()
	dbHost.Start()

	httpHost := http.NewHTTP(fmt.Sprintf(":%d", conf.GlobalConf.Port), dbHost)
	httpHost.Start()

	<-httpHost.Started

	fmt.Println(art)
	fmt.Println(logger.Grey, "	vizn3r's cloud thingy server", logger.Reset)
	fmt.Println(logger.Grey, "	𝘱𝘰𝘸𝘦𝘳𝘦𝘥 𝘣𝘺 𝘨𝘰𝘧𝘪𝘣𝘦𝘳", logger.Reset)
	fmt.Println(logger.Grey, "	© by Simon \"𝘷𝘪𝘻𝘯3𝘳\" Vizner - 2025", logger.Reset)
	fmt.Println()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	fmt.Println()

	log.Close()
	dbHost.Stop()
	httpHost.Stop()

	fmt.Println("Exited gracefully")
}
