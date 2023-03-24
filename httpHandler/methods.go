package httphandler

import "net/http"

func Get(path string, fn http.HandlerFunc) {
	http.HandleFunc(path, fn)
}