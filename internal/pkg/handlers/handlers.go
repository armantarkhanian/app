// Package handlers ...
package handlers

import (
	"app/internal/pkg/server"
	"net/http"

	"github.com/gin-gonic/gin"
)

func index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}
func Init() {
	server.Router.GET("/", index)
}
