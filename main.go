package main

import (
	"fmt"
	"log"

	"app/internal/pkg/websocket"

	"github.com/gin-gonic/gin"
)

func main() {
	websocketHandler, sockJSHandler, err := websocket.RunNode("localhost:6379")
	if err != nil {
		panic(err)
	}

	gin.SetMode(gin.DebugMode)
	router := gin.New()
	router.Use(websocket.GinContextToContextMiddleware())
	router.GET("/connection/websocket", websocketHandler)
	router.GET("/connection/sockjs", sockJSHandler)

	router.LoadHTMLGlob("./*.html")
	router.StaticFile("/index.js", "./index.js")
	router.Static("/static", "./web/static")

	router.GET("/:id", func(c *gin.Context) {
		c.HTML(200, fmt.Sprintf("index%s.html", c.Param("id")), nil)
	})

	router.GET("/check", func(c *gin.Context) {
		online, err := websocket.IsUserOnline("user_15")
		c.JSON(200, gin.H{
			"online": online,
			"error":  err,
		})
	})

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
