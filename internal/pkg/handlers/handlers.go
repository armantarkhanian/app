// Package handlers ...
package handlers

import (
	"app/internal/pkg/jwt"
	"app/internal/pkg/middlewares"
	"app/internal/pkg/recaptcha"
	"app/internal/pkg/server"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

type apiRequest struct {
	Method        string `json:"method"`
	MethodVersion string `json:"methodVersion"`
	Data          string `json:"data"`
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
	server.Router.GET("/", middlewares.Auth(), func(c *gin.Context) {
		userID, tokenID := jwt.GetPayload(c)
		c.JSON(200, gin.H{
			"userID":  userID,
			"tokenID": tokenID,
		})
	})

	server.Router.GET("/login", jwt.LoginHandler())

	server.Router.POST("/internal/api", middlewares.Auth(), func(c *gin.Context) {
		var request apiRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			log.Println("[ERROR]", err)
			c.JSON(200, gin.H{"error": err})
			return
		}
		supportedVersions, exists := supportedMethods[request.Method]
		if !exists {
			c.JSON(200, gin.H{"error": "unsupported method"})
			return
		}
		handlerFunc, exists := supportedVersions[request.MethodVersion]
		if !exists {
			c.JSON(200, gin.H{"error": "unsupported version of method"})
			return
		}
		c.Set("requestData", request.Data)
		handlerFunc(c)
	})
}
