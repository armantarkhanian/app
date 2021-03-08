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

func panicH(c *gin.Context) {
	array := []int{1, 2, 3}
	if len(array) > 5 {
		fmt.Println(array[5])
	}
	c.HTML(http.StatusOK, "index.html", nil)
}

func Init() {
	server.Router.GET("/", index)
	server.Router.GET("/panic", panicH)
}
