package main

import (
	"github.com/gin-gonic/gin"

	"github.com/madeindra/golang-websocket/handler"
)

func main() {
	router := gin.Default()
	router.GET("/socket", handler.WebsocketHandler)
	router.Run(":8080")
}
