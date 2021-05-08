// Package websocket ...
package websocket

import (
	"log"
	"net/http"
	"fmt"
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/centrifugal/centrifuge"
)

func handleLog(e centrifuge.LogEntry) {
	log.Printf("%s: %v", e.Message, e.Fields)
}


type contextKey int

var ginContextKey contextKey

// GinContextToContextMiddleware - at the resolver level we only have access
// to context.Context inside centrifuge, but we need the gin context. So we
// create a gin middleware to add its context to the context.Context used by
// centrifuge websocket server.
func GinContextToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), ginContextKey, c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// GinContextFromContext - we recover the gin context from the context.Context
// struct where we added it just above
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

func RunNode(redisHosts ...string) (gin.HandlerFunc, gin.HandlerFunc, error){
	cfg := centrifuge.DefaultConfig
	cfg.LogLevel = centrifuge.LogLevelInfo
	cfg.LogHandler = handleLog

	node, err := centrifuge.New(cfg)
	if err != nil {
		return nil, nil, err
	}

	if len(redisHosts) > 0 {
		var redisShards []*centrifuge.RedisShard
		for i := 0; i < len(redisHosts); i++ {
			redisShardConfigs := []centrifuge.RedisShardConfig{{Address: redisHosts[i]}}

			for _, redisConf := range redisShardConfigs {
				redisShard, err := centrifuge.NewRedisShard(node, redisConf)
				if err != nil {
					return nil, nil, err
				}
				redisShards = append(redisShards, redisShard)
			}
		}
		broker, err := centrifuge.NewRedisBroker(node, centrifuge.RedisBrokerConfig{
			Shards: redisShards,
		})
		if err != nil {
			return nil, nil, err
		}
		presenceManager, err := centrifuge.NewRedisPresenceManager(node, centrifuge.RedisPresenceManagerConfig{
			Shards: redisShards,
		})
		if err != nil {
		return nil, nil, err
		}
		node.SetBroker(broker)
		node.SetPresenceManager(presenceManager)
	}

	setHandlers(node)

	err = node.Run()
	if err != nil {
		return nil, nil, err
	}

	wsHandler := gin.WrapH(authMiddleware(centrifuge.NewWebsocketHandler(node, centrifuge.WebsocketConfig{
		ReadBufferSize: 1024,
		UseWriteBufferPool: true,
	})))

	sockJSHandler := gin.WrapH(authMiddleware(centrifuge.NewSockjsHandler(node, centrifuge.SockjsConfig{
		URL:           "https://cdn.jsdelivr.net/npm/sockjs-client@1/dist/sockjs.min.js",
		HandlerPrefix: "/connection/sockjs",
		WebsocketReadBufferSize:  1024,
		WebsocketWriteBufferSize: 1024,
	})))
	
	return wsHandler, sockJSHandler, nil
}

func authMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		c, err := GinContextFromContext(ctx) // getGinContext
		if err != nil {
			w.Write([]byte("error, sorry"))			
			return
		}
		username, _ := c.Cookie("user_id")
		if strings.TrimSpace(username) == "" {
			w.Write([]byte("set user_id in cookie (it must be user_15)"))
			return			
		}
		newCtx := centrifuge.SetCredentials(ctx, &centrifuge.Credentials{
			UserID: username,
		})
		r = r.WithContext(newCtx)
		h.ServeHTTP(w, r)
	})
}
