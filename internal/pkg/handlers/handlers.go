// Package handlers ...
package handlers

import (
	"app/internal/pkg/middlewares"
	"app/internal/pkg/recaptcha"
	"app/internal/pkg/server"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type apiRequest struct {
	Data string `json:"data"`
}

var (
	supportedMethods = map[string]map[string]gin.HandlerFunc{
		"createAccount": {
			"v1": createAccount1,
		},
		"verifyCaptcha": {
			"v1": recaptcha.VerifyCaptcha1,
		},
	}
)

func Init() {
	server.Router.GET("/panic", func(c *gin.Context) {
		array := []int{1, 2, 3}
		fmt.Println(array[5])
		c.String(200, "Hello, World")
	})
	server.Router.GET("/", func(c *gin.Context) {
		time.Sleep(33 * time.Millisecond)
		c.String(200, "Hello, World")
	})
	server.Router.POST("/internal/login", login)
	server.Router.POST("/internal/logout", logout)

	server.Router.POST("/internal/api/:m/:v", middlewares.Auth(), middlewares.RecaptchaProtected(), func(c *gin.Context) {
		method := c.Param("m")
		version := c.Param("v")
		var request apiRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			log.Println("[ERROR]", err)
			c.JSON(200, gin.H{"error": err})
			return
		}
		supportedVersions, exists := supportedMethods[method]
		if !exists {
			c.JSON(200, gin.H{"error": "unsupported method"})
			return
		}
		handlerFunc, exists := supportedVersions[version]
		if !exists {
			c.JSON(200, gin.H{"error": "unsupported version of method"})
			return
		}
		c.Set("requestData", request.Data)
		handlerFunc(c)
	})
}
