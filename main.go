package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()
	router.Use(GinContextToContextMiddleware())

	router.GET("/connection/websocket", centrifugeHandler())
	router.LoadHTMLFiles("index.html")

	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
