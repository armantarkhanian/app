// Package websocket ...
package websocket

import (
	"log"
	"net/http"

	"github.com/centrifugal/centrifuge"
)

func handleLog(e centrifuge.LogEntry) {
	log.Printf("%s: %v", e.Message, e.Fields)
}

func Run(addr string, redisHosts ...string) {
	cfg := centrifuge.DefaultConfig
	cfg.LogLevel = centrifuge.LogLevelInfo
	cfg.LogHandler = handleLog

	node, err := centrifuge.New(cfg)
	if err != nil {
		log.Fatalln(err)
	}

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
			log.Fatalln(err)
		}
		presenceManager, err := centrifuge.NewRedisPresenceManager(node, centrifuge.RedisPresenceManagerConfig{
			Shards: redisShards,
		})
		if err != nil {
			log.Fatalln(err)
		}
		node.SetBroker(broker)
		node.SetPresenceManager(presenceManager)
	}

	setHandlers(node)

	err = node.Run()
	if err != nil {
		log.Fatalln(err)
	}
	
	http.Handle("/connection/websocket", authMiddleware(centrifuge.NewWebsocketHandler(node, centrifuge.WebsocketConfig{
		ReadBufferSize: 1024,
		CheckOrigin: func(*http.Request) bool {
			return true
		},
		UseWriteBufferPool: true,
	})))
	http.Handle("/connection/sockjs/", authMiddleware(centrifuge.NewSockjsHandler(node, centrifuge.SockjsConfig{
		URL:           "https://cdn.jsdelivr.net/npm/sockjs-client@1/dist/sockjs.min.js",
		HandlerPrefix: "/connection/sockjs",
		CheckOrigin: func(*http.Request) bool {
			return true
		},
		WebsocketReadBufferSize:  1024,
		WebsocketWriteBufferSize: 1024,
	})))

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}

func authMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		newCtx := centrifuge.SetCredentials(ctx, &centrifuge.Credentials{
			UserID: "user_15",
		})
		r = r.WithContext(newCtx)
		h.ServeHTTP(w, r)
	})
}
