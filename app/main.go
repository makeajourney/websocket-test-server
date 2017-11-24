package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

const (
	socketBufferSize = 1024
)

var (
	renderer *render.Render
	upgrader = &websocket.Upgrader{
		ReadBufferSize:  socketBufferSize,
		WriteBufferSize: socketBufferSize,
	}
)

func init() {
	renderer = render.New()
}

func main() {
	router := httprouter.New()

	router.GET("/", func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		renderer.HTML(w, http.StatusOK, "index", nil)
	})

	router.GET("/ws/", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		socket, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal("ServeHTTP:", err)
			return
		}
		newClient(socket)
	})

	n := negroni.Classic()
	n.UseHandler(router)
	n.Run(":3000")
}
