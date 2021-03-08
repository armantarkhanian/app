package main

import (
	"app/internal/pkg/configs"
	"app/internal/pkg/handlers"
	"app/internal/pkg/server"
	"app/internal/pkg/sessions"
	"log"
)

func init() {
	configs.Init()
	sessions.Init()
	server.Init()
	handlers.Init()
}

func main() {
	log.Println("This is a test log entry")

	server.Run()
}
