// Package redis ...
package redis

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

var (
	Ctx    = context.Background()
	Client *redis.Client
)

func Init() {
	Client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	if err := Client.Ping(Ctx).Err(); err != nil {
		log.Fatalln("[FATAL] [redis]", err)
	}
}
