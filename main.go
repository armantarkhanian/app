package main

import (
	"app/internal/handlers"
	"app/internal/pkg/configs"
	"app/internal/pkg/server"
	"app/internal/pkg/sessions"
)

func init() {
	configs.Init()
	sessions.Init()
	server.Init()
	handlers.Init()
}

func main() {
	server.Run()
}
