// Package websocket ...
package websocket

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/centrifugal/centrifuge"
)

func Run(addr string, redisHosts ...string) error {
	node, err := centrifuge.New(centrifuge.DefaultConfig)
	if err != nil {
		return err
	}
	setHandlers(node)

	if len(redisHosts) > 0 {
		var redisShards []*centrifuge.RedisShard
		for i := 0; i < len(redisHosts); i++ {
			redisShardConfigs := []centrifuge.RedisShardConfig{{Address: redisHosts[i]}}

			for _, redisConf := range redisShardConfigs {
				redisShard, err := centrifuge.NewRedisShard(node, redisConf)
				if err != nil {
					log.Println(err)
				}
				redisShards = append(redisShards, redisShard)
			}
		}
		broker, err := centrifuge.NewRedisBroker(node, centrifuge.RedisBrokerConfig{
			Shards: redisShards,
		})
		if err != nil {
			return err
		}
		presenceManager, err := centrifuge.NewRedisPresenceManager(node, centrifuge.RedisPresenceManagerConfig{
			Shards: redisShards,
		})
		if err != nil {
			return err
		}
		node.SetBroker(broker)
		node.SetPresenceManager(presenceManager)
	}

	err = node.Run()
	if err != nil {
		return err
	}

	websocketHandler := centrifuge.NewWebsocketHandler(node, centrifuge.WebsocketConfig{
		ReadBufferSize:     1024,
		UseWriteBufferPool: true,
	})

	http.Handle("/connection/websocket", authMiddleware(websocketHandler))

	sockjsHandler := centrifuge.NewSockjsHandler(node, centrifuge.SockjsConfig{
		URL:                      "https://cdn.jsdelivr.net/npm/sockjs-client@1/dist/sockjs.min.js",
		HandlerPrefix:            "/connection/sockjs",
		WebsocketReadBufferSize:  1024,
		WebsocketWriteBufferSize: 1024,
	})
	http.Handle("/connection/sockjs/", authMiddleware(sockjsHandler))

	http.Handle("/", http.FileServer(http.Dir("./")))

	return nil
}

func authMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// r.Header.Get("Authorization")
		// or c.Cookie("jwt")
		// extract claims and set userID to context
		ctx := r.Context()
		newCtx := centrifuge.SetCredentials(ctx, &centrifuge.Credentials{
			UserID: "admin",
		})
		r = r.WithContext(newCtx)
		h.ServeHTTP(w, r)
	})
}

func waitExitSignal(n *centrifuge.Node) {
	sigCh := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		_ = n.Shutdown(context.Background())
		done <- true
	}()
	<-done
}
