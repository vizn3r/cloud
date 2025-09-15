package main

import (
	"cloud-server/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	httpHost := http.NewHTTP(":8080")
	httpHost.Start()
	defer httpHost.Stop()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}
