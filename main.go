package main

import (
	_ "app/internal/pkg/init"
	"app/internal/pkg/server"
)

func main() {
	server.Run()
}
