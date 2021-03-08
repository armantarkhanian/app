// Package handlers ...
package handlers

import (
	"app/internal/pkg/recaptcha"
	"app/internal/pkg/server"
	"time"

	"github.com/gin-contrib/sessions"

	"github.com/gin-gonic/gin"
)

func ok(c *gin.Context) {
	session := sessions.Default(c)
	session.Set("needCaptcha", false)
	session.Set("qps", 0)
	session.Set("lastActionTime", time.Time{})
	session.Save()
}

func index(c *gin.Context) {
	session := sessions.Default(c)
	lastActionTime, _ := session.Get("lastActionTime").(time.Time)
	qps, _ := session.Get("qps").(int)

	c.IndentedJSON(200, gin.H{
		"lastActionTime": lastActionTime,
		"qps":            qps,
	})
}

func Init() {
	server.Router.GET("/", recaptcha.Middleware(), index)
	server.Router.GET("/ok", ok)
}
