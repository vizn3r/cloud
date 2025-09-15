package main

import (
	"cloud-server/conf"
	"cloud-server/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	conf.LoadConfig("./server.json")

	httpHost := http.NewHTTP(":8080")
	httpHost.Start()
	defer httpHost.Stop()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}
