// Package handlers ...
package handlers

import (
	"app/internal/pkg/configs"
	"app/internal/pkg/metrics"
	"app/internal/pkg/server"
	"app/internal/pkg/sessions"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func femaleHandler(c *gin.Context) {
	session := sessions.Get(c)
	session.UserID = "admin"
	session.Gender = "female"
	session.Save(c)
}

func maleGender(c *gin.Context) {
	session := sessions.Get(c)
	session.UserID = "admin"
	session.Gender = "male"
	session.Save(c)
}

func panicHandler(c *gin.Context) {
	time.Sleep(1 * time.Second)
	c.IndentedJSON(http.StatusOK, configs.Store)
}

func Init() {
	server.Router.GET("/female", femaleHandler)
	server.Router.GET("/male", maleGender)
	server.Router.GET("/metrics", metrics.InfoHandler())
}
