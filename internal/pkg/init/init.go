// Package init ...
package init

import (
	"app/internal/pkg/configs"
	"app/internal/pkg/db"
	"app/internal/pkg/handlers"
	"app/internal/pkg/logger"
	"app/internal/pkg/redis"
	"app/internal/pkg/server"
	"app/internal/pkg/sessions"
)

func init() {
	configs.Init()
	if err := logger.Init(); err != nil {
		log.Fatalln("[FATAL] [logger]", err)
	}
	
	sessions.Init()
	db.Init()
	redis.Init()
	server.Init()
	handlers.Init()
}
