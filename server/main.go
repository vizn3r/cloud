package main

import (
	"cloud-server/conf"
	"cloud-server/db"
	"cloud-server/http"
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

//go:embed server.json
var config string

func main() {
	if err := conf.LoadFromBytes([]byte(config)); err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat("storage"); os.IsNotExist(err) {
		err := os.MkdirAll("storage/temp", 0700)
		if err != nil {
			log.Fatal(err)
		}
	}

	dbHost := db.NewDB()
	dbHost.Start()
	defer dbHost.Stop()

	httpHost := http.NewHTTP(fmt.Sprintf(":%d", conf.GlobalConf.Port), dbHost)
	httpHost.Start()
	defer httpHost.Stop()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}
