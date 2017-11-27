package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	socketBufferSize = 1024
)

var (
	router   *gin.Engine
	upgrader = &websocket.Upgrader{
		ReadBufferSize:  socketBufferSize,
		WriteBufferSize: socketBufferSize,
	}
)

func main() {
	router = gin.Default()

	router.LoadHTMLGlob("templates/*")

	initializeRoutes()

	router.Run()
}
