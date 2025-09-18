package main

import (
	"cloud-server/conf"
	"cloud-server/db"
	"cloud-server/http"
	"cloud-server/logger"
	_ "embed"
	"fmt"
	"os"
	"os/signal"
	"syscall"
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

	// Ensure temp directory exists
	if _, err := os.Stat("storage/temp"); os.IsNotExist(err) {
		err := os.MkdirAll("storage/temp", 0700)
		if err != nil {
			log.Fatal(err)
		}
	}

	dbHost := db.NewDB()
	dbHost.Start()

	httpHost := http.NewHTTP(fmt.Sprintf(":%d", conf.GlobalConf.Port), dbHost)
	httpHost.Start()

	log.Print(art)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	fmt.Println()

	log.Close()
	dbHost.Stop()
	httpHost.Stop()

	fmt.Println("Exited gracefully")
}
