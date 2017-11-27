package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func showIndexPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", nil)
}

func wsConnection(c *gin.Context) {
	socket, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatal("ServeHTTP: ", err)
		return
	}
	newClient(socket)
}
