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
	p, _ := os.Getwd()
	fmt.Println(p)
	godotenv.Load( p + "/.env")
	uploader.Start(PORT)
	fmt.Println("Started at port" + PORT)
	http.ListenAndServe(PORT, nil)
}