package main

import (
	"cloud/uploader"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	uploader.Start()
	http.ListenAndServe(":9873", nil)
}