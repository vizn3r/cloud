package main

import (
	"cloud/test"
)

const (
	PORT = ":9873"
)

func main() {
	// godotenv.Load(".env")
	// uploader.Start(PORT)
	// fmt.Println("Started at port" + PORT)
	// http.ListenAndServe(PORT, nil)
	test.Test()
}