package client

import (
	h "cloud/handler"
	"net/http"
)

func Start(port int) {
	http.HandleFunc("/ui", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, h.Cwd("/"))	
	})
}