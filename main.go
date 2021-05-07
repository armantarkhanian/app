package main

import (
	"fmt"
	"log"

	"github.com/centrifugal/centrifuge"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()
	router.GET("/connection/websocket", GinContextToContextMiddleware(), func(c *gin.Context) {
		username, _ := c.Cookie("username")
		fmt.Println(username)
		cred := &centrifuge.Credentials{
			UserID: username,
		}
		newCtx := centrifuge.SetCredentials(c.Request.Context(), cred)
		c.Request = c.Request.WithContext(newCtx)
		c.Next()
	}, centrifugeHandler())
	router.LoadHTMLFiles("index.html")
	router.StaticFile("/index.js", "./index.js")
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
