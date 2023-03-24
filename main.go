package main

import (
	"cloud/uploader"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

const (
	PORT = ":9873"
)

func main() {
	godotenv.Load(".env")
	fmt.Println(os.Getenv("AUTHORIZATION"))
	uploader.Start(PORT)
	fmt.Println("Started at port" + PORT)
	http.ListenAndServe(PORT, nil)
}