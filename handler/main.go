package handler

import (
	"fmt"
	"net/http"
)

type handler struct {
	routers []router
	port    int
}

func HandlerInit(port int) handler {
	return handler{nil, port}
}

func (h *handler) Start() {
	println("STARTING HANDLER")
	for _, r := range h.routers {
		r.Load()
	}
	print("LISTENING ON PORT: ")
	println(h.port)
	http.ListenAndServe(":" + fmt.Sprint(h.port), nil)
}