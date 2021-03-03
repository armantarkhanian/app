// Package server ...
package server

import (
	"app/internal/pkg/configs"
	"app/internal/pkg/geoip"
	"app/internal/pkg/logger"
	"app/internal/pkg/sessions"
	"io"
	"net/http"
	"os"
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

	if configs.Store.Gin.Log.AccessLogFile != "" {
		accessLogFile, err := os.OpenFile(configs.Store.Gin.Log.AccessLogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalln(err)
		}
		if configs.Store.Gin.Log.UseStdOut {
			gin.DefaultWriter = io.MultiWriter(os.Stdout, accessLogFile)
		} else {
			gin.DefaultWriter = io.MultiWriter(accessLogFile)
		}
	}
	if configs.Store.Gin.Log.ErrorLogFile != "" {
		errorLogFile, err := os.OpenFile(configs.Store.Gin.Log.ErrorLogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalln(err)
		}
		if configs.Store.Gin.Log.UseStdErr {
			gin.DefaultErrorWriter = io.MultiWriter(os.Stderr, errorLogFile)
		} else {
			gin.DefaultErrorWriter = io.MultiWriter(errorLogFile)
		}
	}

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
	Router.Use(logger.Middleware())
	Router.Use(gin.Recovery())
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
