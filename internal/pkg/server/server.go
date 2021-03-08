// Package server ...
package server

import (
	"app/internal/pkg/configs"
	"app/internal/pkg/geoip"
	"app/internal/pkg/logger"
	"app/internal/pkg/sessions"
	"context"
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

	shutdownTimeout, err = time.ParseDuration(configs.Store.Gin.Timeouts.Shutdown)
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

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("[INFO] Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("[INFO] Server forced to shutdown:", err)
	}

	log.Println("[INFO] Server exiting")

}
