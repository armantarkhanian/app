// Package server ...
package server

import (
	"app/internal/pkg/configs"
	"app/internal/pkg/middlewares"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"log"

	"github.com/gin-gonic/gin"
)

var (
	Router *gin.Engine
	server *http.Server

	readTimeout     time.Duration
	writeTimeout    time.Duration
	idleTimeout     time.Duration
	shutdownTimeout time.Duration
)

func Init() {
	gin.SetMode(configs.Store.Gin.Mode)

	readTimeout, err := time.ParseDuration(configs.Store.Gin.Timeouts.Read)
	if err != nil {
		log.Fatalln("[FATAL]", err)
	}
	writeTimeout, err = time.ParseDuration(configs.Store.Gin.Timeouts.Write)
	if err != nil {
		log.Fatalln("[FATAL]", err)
	}
	idleTimeout, err = time.ParseDuration(configs.Store.Gin.Timeouts.Idle)
	if err != nil {
		log.Fatalln("[FATAL]", err)
	}
	shutdownTimeout, err = time.ParseDuration(configs.Store.Gin.Timeouts.Shutdown)
	if err != nil {
		log.Fatalln("[FATAL]", err)
	}

	Router = gin.New()
	Router.Use(middlewares.Recovery())
	Router.Use(middlewares.GeoIP())
	Router.Use(middlewares.Sessions())
	Router.Use(middlewares.AccessLogger())

	Router.Delims("[[", "]]")
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

	go func() {
		log.Println(fmt.Sprintf("[INFO] running server on port %q", configs.Store.Gin.Addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln("[FATAL]", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("[INFO] shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalln("[WARNING] server forced to shutdown:", err)
	}

	log.Println("[INFO] server exiting properly")
}
