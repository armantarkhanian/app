package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	// Import this library.
	"github.com/centrifugal/centrifuge"
	"github.com/gin-gonic/gin"
)

type contextKey int

var ginContextKey contextKey

func GinContextToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), ginContextKey, c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func GinContextFromContext(ctx context.Context) (*gin.Context, error) {
	ginContext := ctx.Value(ginContextKey)
	if ginContext == nil {
		err := fmt.Errorf("could not retrieve gin.Context")
		return nil, err
	}
	gc, ok := ginContext.(*gin.Context)
	if !ok {
		err := fmt.Errorf("gin.Context has wrong type")
		return nil, err
	}
	return gc, nil
}

func authMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		c, err := GinContextFromContext(ctx)
		if err != nil {
			fmt.Printf("Failed to retrieve gin context")
			fmt.Print(err.Error())
			return
		}
		username, _ := c.Cookie("username")
		if username != "user_15" {
			w.Write([]byte(`{"error": "unauthorised"}`))
			return
		}
		cred := &centrifuge.Credentials{
			UserID: username,
		}
		newCtx := centrifuge.SetCredentials(ctx, cred)
		r = r.WithContext(newCtx)
		h.ServeHTTP(w, r)
	})
}

func initNode() (*centrifuge.Node, error) {
	node, err := centrifuge.New(centrifuge.DefaultConfig)
	if err != nil {
		return nil, err
	}

	connectHandler := func(client *centrifuge.Client) {
		client.OnRPC(func(e centrifuge.RPCEvent, callback centrifuge.RPCCallback) {
			callback(onRPC(node, client, &e))
		})
		client.OnSubscribe(func(e centrifuge.SubscribeEvent, callback centrifuge.SubscribeCallback) {
			callback(onSubscribe(node, client, &e))
		})
		client.OnMessage(func(e centrifuge.MessageEvent) {
			onMessage(node, client, &e)
		})
		client.OnPublish(func(e centrifuge.PublishEvent, callback centrifuge.PublishCallback) {
			callback(onPublish(node, client, &e))
		})
		client.OnDisconnect(func(e centrifuge.DisconnectEvent) {
			onDisconect(node, client, &e)
		})
	}

	node.OnConnect(connectHandler)

	redisShardConfigs := []centrifuge.RedisShardConfig{{Address: "localhost:6379"}}

	var redisShards []*centrifuge.RedisShard
	for _, redisConf := range redisShardConfigs {
		redisShard, err := centrifuge.NewRedisShard(node, redisConf)
		if err != nil {
			log.Println(err)
		}
		redisShards = append(redisShards, redisShard)
	}

	broker, err := centrifuge.NewRedisBroker(node, centrifuge.RedisBrokerConfig{
		Shards: redisShards,
	})
	if err != nil {
		log.Println(err)
	}

	node.SetBroker(broker)

	presenceManager, err := centrifuge.NewRedisPresenceManager(node, centrifuge.RedisPresenceManagerConfig{
		Shards: redisShards,
	})
	if err != nil {
		log.Println(err)
	}
	node.SetPresenceManager(presenceManager)

	return node, node.Run()
}

func centrifugeHandler() gin.HandlerFunc {
	node, err := initNode()
	if err != nil {
		panic(err)
	}
	wsHandler := centrifuge.NewWebsocketHandler(node, centrifuge.WebsocketConfig{})
	return gin.WrapH(wsHandler)
}
