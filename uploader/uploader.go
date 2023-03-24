package uploader

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

var mimeTypes = map[string]string{
	"image/png": ".png",
	"image/gif": ".gif",
	"video/mp4": ".mp4",
}

func Start(PORT string) {
	rand.Seed(time.Now().UnixMicro())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if r.Method == "GET" {
			if strings.Contains(path, "..") {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			if path == "/" {
				w.WriteHeader(http.StatusBadRequest)
			}
			http.ServeFile(w, r, "./files" + path)
		} else if r.Method == "POST" {
			if auth := r.Header.Get("authorization"); auth != os.Getenv("AUTHORIZATION") {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			defer r.Body.Close()
			data, _ := ioutil.ReadAll(r.Body)
			rand := fmt.Sprintf("%x", rand.Int63())
			ext := mimeTypes[r.Header.Get("content-type")]
			ioutil.WriteFile("./files/" + rand + ext, data, 0644)
			w.Write([]byte("http://localhost:" + PORT + "/" + rand + ext))
		}
	})	
}