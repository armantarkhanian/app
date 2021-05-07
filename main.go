package main

import (
	"fmt"
	"log"

	"app/internal/pkg/websocket"

	"github.com/gin-gonic/gin"
)

func main() {
	go func() {
		websocket.Run(":8000", "localhost:6379")
	}()

	router := gin.New()

	router.LoadHTMLFiles("index.html")
	router.StaticFile("/index.js", "./index.js")

	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
