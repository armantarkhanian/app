// Package server ...
package server

import (
	"app/internal/pkg/configs"
	"app/internal/pkg/geoip"
	"app/internal/pkg/logger"
	"app/internal/pkg/sessions"
	"net/http"
	"time"

	"log"

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

	readTimeout, err := time.ParseDuration(configs.Store.Gin.Timeouts.Read)
	if err != nil {
		log.Fatalln(err)
	}
	writeTimeout, err = time.ParseDuration(configs.Store.Gin.Timeouts.Write)
	if err != nil {
		log.Fatalln(err)
	}
	idleTimeout, err = time.ParseDuration(configs.Store.Gin.Timeouts.Idle)
	if err != nil {
		log.Fatalln(err)
	}

	Router = gin.New()
	Router.Use(logger.Recovery())
	Router.Use(geoip.Middleware())
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
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
