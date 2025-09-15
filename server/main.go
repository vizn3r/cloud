package main

import (
	"cloud-server/conf"
	"cloud-server/db"
	"cloud-server/http"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	conf.LoadConfig("./server.json")

	if _, err := os.Stat("storage"); os.IsNotExist(err) {
		err := os.MkdirAll("storage/temp", 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	dbHost := db.NewDB()
	dbHost.Start()
	defer dbHost.Stop()

	httpHost := http.NewHTTP(fmt.Sprintf(":%d", conf.GlobalConf.Port))
	httpHost.Start()
	defer httpHost.Stop()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}
