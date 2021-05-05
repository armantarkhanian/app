// Package redis ...
package redis

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

var (
	Ctx     = context.Background()
	Cluster *redis.ClusterClient
)

func Init() {
	Cluster = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{"localhost:6379"},
	})
	if err := Cluster.Ping(Ctx).Err(); err != nil {
		log.Fatalln("[FATAL] [redis]", err)
	}
}
