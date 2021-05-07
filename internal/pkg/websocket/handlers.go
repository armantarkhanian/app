// Package websocket ...
package websocket

import (
	"app/internal/pkg/handlers"
	"context"

	"github.com/centrifugal/centrifuge"
)

func setHandlers(node *centrifuge.Node) {
	// first connection, subscirbe user to his own channel
	node.OnConnecting(func(ctx context.Context, e centrifuge.ConnectEvent) (centrifuge.ConnectReply, error) {
		cred, _ := centrifuge.GetCredentials(ctx)
		return centrifuge.ConnectReply{
			Data: []byte(`{}`),
			// Subscribe to personal several server-side channel.
			Subscriptions: map[string]centrifuge.SubscribeOptions{
				"#" + cred.UserID: {Recover: true, Presence: true, JoinLeave: true},
			},
		}, nil
	})

	node.OnConnect(func(client *centrifuge.Client) {
		client.OnAlive(func() {
			handlers.AliveHandler(node, client)
		})
		client.OnDisconnect(func(e centrifuge.DisconnectEvent) {
			handlers.DisconectHandler(node, client, &e)
		})
		client.OnSubscribe(func(e centrifuge.SubscribeEvent, callback centrifuge.SubscribeCallback) {
			callback(handlers.SubscribeHandler(node, client, &e))
		})
		client.OnUnsubscribe(func(e centrifuge.UnsubscribeEvent) {
			handlers.UnsubscribeHandler(node, client, &e)
		})
		client.OnPublish(func(e centrifuge.PublishEvent, callback centrifuge.PublishCallback) {
			callback(handlers.PublishHandler(node, client, &e))
		})
		client.OnRefresh(func(e centrifuge.RefreshEvent, callback centrifuge.RefreshCallback) {
			callback(handlers.RefreshHandler(node, client, &e))
		})
		client.OnSubRefresh(func(e centrifuge.SubRefreshEvent, callback centrifuge.SubRefreshCallback) {
			callback(handlers.SubRefreshHandler(node, client, &e))
		})
		client.OnRPC(func(e centrifuge.RPCEvent, callback centrifuge.RPCCallback) {
			callback(handlers.RPCHandler(node, client, &e))
		})
		client.OnMessage(func(e centrifuge.MessageEvent) {
			handlers.MessageHandler(node, client, &e)
		})
		client.OnPresence(func(e centrifuge.PresenceEvent, callback centrifuge.PresenceCallback) {
			callback(handlers.PresenceHandler(node, client, &e))
		})
		client.OnPresenceStats(func(e centrifuge.PresenceStatsEvent, callback centrifuge.PresenceStatsCallback) {
			callback(handlers.PresenceStatsHandler(node, client, &e))
		})
		client.OnHistory(func(e centrifuge.HistoryEvent, callback centrifuge.HistoryCallback) {
			callback(handlers.HistoryHandler(node, client, &e))
		})
	})
}
