package main

import (
	"cloud/uploader"
	"fmt"
	"net/http"

	"github.com/joho/godotenv"
)

const (
	PORT = ":9873"
)

func main() {
	godotenv.Load(".env")
	uploader.Start(PORT)
	fmt.Println("Started at port" + PORT)
	http.ListenAndServe(PORT, nil)
}