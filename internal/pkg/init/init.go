// Package init ...
package init

import (
	"app/internal/pkg/configs"
	"app/internal/pkg/handlers"
	"app/internal/pkg/logger"
	"app/internal/pkg/redis"
	"app/internal/pkg/server"
	"app/internal/pkg/sessions"
)

func init() {
	configs.Init()
	logger.Init()
	sessions.Init()
	redis.Init()
	server.Init()
	handlers.Init()
}
