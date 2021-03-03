// Package handlers ...
package handlers

import (
	"app/internal/pkg/server"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func index(c *gin.Context) {
	v, _ := c.Get("countryCode")
	countryCode, _ := v.(string)
	fmt.Println(countryCode)
	c.HTML(http.StatusOK, "index.html", nil)
}

func Init() {
	server.Router.GET("/", index)
}
