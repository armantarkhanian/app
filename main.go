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
	router.Static("/static", "./web/static/")
	router.StaticFile("/index.js", "./index.js")

	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	type privateChannel struct {
		UserID string `json:"client"`
		Channels []string `json:"channels"`
	}
	router.POST("/centrifuge/subscribe", func(c *gin.Context) {
		var e privateChannel
		if err := c.ShouldBindJSON(&e); err != nil {
			c.JSON(200, gin.H{
				"channels": []string{},
			})
			return
		}
		fmt.Println(e)		
		c.JSON(200, gin.H{
			"channels": []string{},
		})
	})

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
