// Package handlers ...
package handlers

import (
	"app/internal/pkg/middlewares"
	"app/internal/pkg/recaptcha"
	"app/internal/pkg/server"
	"log"

	"github.com/gin-gonic/gin"
)

type apiRequest struct {
	Method  string `json:"method"`
	Version string `json:"version"`
	Data    string `json:"data"`
}

var (
	supportedMethods = map[string]map[string]gin.HandlerFunc{
		"createAccount": {
			"1.0": createAccount1,
		},
		"verifyCaptcha": {
			"1.0": recaptcha.VerifyCaptcha,
		},
	}
)

func Init() {
	server.Router.POST("/internal/api", middlewares.RecaptchaProtect(), func(c *gin.Context) {
		var request apiRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			log.Println("[ERROR]", err)
			c.JSON(200, gin.H{"error": err})
			return
		}
		handlerFunc, exists := supportedMethods[request.Method][request.Version]
		if !exists {
			c.JSON(200, gin.H{"error": "unsupported method or version"})
			return
		}
		c.Set("requestData", request.Data)
		handlerFunc(c)
	})
}
