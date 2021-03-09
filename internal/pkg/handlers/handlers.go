// Package handlers ...
package handlers

import (
	"app/internal/pkg/recaptcha"
	"app/internal/pkg/server"

	"github.com/gin-gonic/gin"
)

func index(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{
		"message": "",
	})
}

func users(c *gin.Context) {
	type user struct {
		Name string `json:"name"`
	}
	type response struct {
		Users []user `json:"users"`
	}
	resp := response{
		Users: []user{
			{"arman 1"},
			{"arman 2"},
			{"arman 3"},
			{"arman 4"},
		},
	}
	c.JSON(200, resp)
}

func Init() {
	server.Router.GET("/", index)
	server.Router.POST("/ok", recaptcha.CheckCaptcha)
	server.Router.GET("/users", recaptcha.Middleware(), users)
}
