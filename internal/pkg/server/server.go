// Package server ...
package server

import (
	"app/internal/pkg/configs"
	"app/internal/pkg/sessions"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	Router *gin.Engine
	server *http.Server

	readTimeout  time.Duration
	writeTimeout time.Duration
	idleTimeout  time.Duration
)

func Init() {
	gin.SetMode(configs.Store.Gin.Mode)

	var err error
	readTimeout, err = time.ParseDuration(configs.Store.Gin.Timeouts.Read)
	if err != nil {
		panic(err)
	}
	writeTimeout, err = time.ParseDuration(configs.Store.Gin.Timeouts.Write)
	if err != nil {
		panic(err)
	}
	idleTimeout, err = time.ParseDuration(configs.Store.Gin.Timeouts.Idle)
	if err != nil {
		panic(err)
	}

	Router = gin.New()
	Router.Use(gin.Recovery())
	Router.Use(sessions.Middleware())
	Router.LoadHTMLGlob("./web/template/*.html")
	Router.Static("/static", "./web/static/")

	server = &http.Server{
		Addr:         configs.Store.Gin.Addr,
		Handler:      Router,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}
}

func Run() {
	log.Fatalln(server.ListenAndServe())
}
