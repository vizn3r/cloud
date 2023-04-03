package uploader

// import (
// 	h "cloud/handler"
// 	"fmt"
// 	"io/ioutil"
// 	"math/rand"
// 	"net/http"
// 	"os"
// 	"strings"
// 	"time"
// )

// var mimeTypes = map[string]string{
// 	"image/png": ".png",
// 	"image/gif": ".gif",
// 	"video/mp4": ".mp4",
// }

// func Start(PORT string) { // replace PORT with config
// 	rand.Seed(time.Now().UnixMicro())
// 	rt := h.Route("/")
// 	rt.Fn("GET", func(w http.ResponseWriter, r *http.Request) {
// 		path := r.URL.Path
// 		if strings.Contains(path, "..") {
// 			w.WriteHeader(http.StatusForbidden)
// 			return
// 		}
// 		if path == "/" {
// 			w.WriteHeader(http.StatusBadRequest)
// 			return
// 		}
// 		http.ServeFile(w, r, h.Cwd("/files") + path)
// 	})
// 	rt.Fn("POST", func(w http.ResponseWriter, r *http.Request) {
// 		if auth := r.Header.Get("authorization"); auth != os.Getenv("AUTHORIZATION") {
// 			w.WriteHeader(http.StatusUnauthorized)
// 			return
// 		}
// 		defer r.Body.Close()
// 		data, _ := ioutil.ReadAll(r.Body)
// 		rand := fmt.Sprintf("%x", rand.Int63())
// 		ext := mimeTypes[r.Header.Get("content-type")]
// 		ioutil.WriteFile(h.Cwd("/files") + rand + ext, data, 0644)
// 		w.Write([]byte("https://" + r.Host + "/" + rand + ext))
// 	})
// }