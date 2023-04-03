package test

import (
	h "cloud/handler"
	"log"
	"os"
)

func Test() {
	hnd := h.HandlerInit(987)
	rtr := h.RouterInit(hnd)
	rtr.All("/", func(req h.Req, res h.Res) {
		req.Body.Close()
		file, err := os.ReadFile("C:/Users/simon/Desktop/code/vizn3r/cloud/client/index.html")
		if err != nil {
				log.Fatal(err)
			}
		res.Write(file)
	})
	hnd.Start()
}