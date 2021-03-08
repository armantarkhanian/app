// Package handlers ...
package handlers

import (
	"app/internal/pkg/server"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func friends(c *gin.Context) {
	b := 0
	var id int = 15 / b
	fmt.Println(id)
	c.HTML(http.StatusOK, "index.html", nil)
}

func Init() {
	server.Router.GET("/", index)
	server.Router.GET("/friends", friends)
}
