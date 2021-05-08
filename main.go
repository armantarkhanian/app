package main

import (
	"log"

	"app/internal/pkg/websocket"

	"github.com/gin-gonic/gin"
)

func main() {
	websocketHandler, sockJSHandler, err := websocket.RunNode("localhost:6379")
	if err != nil {
		panic(err)
	}

	router := gin.New()
	router.Use(websocket.GinContextToContextMiddleware())
	router.GET("/connection/websocket", websocketHandler)
	router.GET("/connection/sockjs", sockJSHandler)

	router.LoadHTMLFiles("index.html")
	router.StaticFile("/index.js", "./index.js")
	router.GET("/", func(c *gin.Context) {c.HTML(200, "index.html", nil)})

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
