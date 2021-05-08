// Package websocket ...
package websocket

import (
	"context"
	"fmt"

	"github.com/centrifugal/centrifuge"
)

func setHandlers(node *centrifuge.Node) {
	node.OnConnecting(func(ctx context.Context, e centrifuge.ConnectEvent) (centrifuge.ConnectReply, error) {
		cred, ok := centrifuge.GetCredentials(ctx)
		fmt.Println(ok, cred.UserID)
		return centrifuge.ConnectReply{
			Subscriptions: map[string]centrifuge.SubscribeOptions{
				"#" + cred.UserID: {
					Recover:   true,
					Presence:  true,
					JoinLeave: true,
				},
			},
		}, nil
	})

	node.OnConnect(func(client *centrifuge.Client) {		
		transportName := client.Transport().Name()
		transportProto := client.Transport().Protocol()
		fmt.Printf("Client %q is connect via %q with %q protocol", client.UserID(), transportName, transportProto)
		client.OnAlive(func() {
			AliveHandler(node, client)
		})
		client.OnDisconnect(func(e centrifuge.DisconnectEvent) {
			DisconectHandler(node, client, &e)
		})
		client.OnSubscribe(func(e centrifuge.SubscribeEvent, callback centrifuge.SubscribeCallback) {
			callback(SubscribeHandler(node, client, &e))
		})
		client.OnUnsubscribe(func(e centrifuge.UnsubscribeEvent) {
			UnsubscribeHandler(node, client, &e)
		})
		client.OnPublish(func(e centrifuge.PublishEvent, callback centrifuge.PublishCallback) {
			callback(PublishHandler(node, client, &e))
		})
		client.OnRefresh(func(e centrifuge.RefreshEvent, callback centrifuge.RefreshCallback) {
			callback(RefreshHandler(node, client, &e))
		})
		client.OnSubRefresh(func(e centrifuge.SubRefreshEvent, callback centrifuge.SubRefreshCallback) {
			callback(SubRefreshHandler(node, client, &e))
		})
		client.OnRPC(func(e centrifuge.RPCEvent, callback centrifuge.RPCCallback) {
			callback(RPCHandler(node, client, &e))
		})
		client.OnMessage(func(e centrifuge.MessageEvent) {
			MessageHandler(node, client, &e)
		})
		client.OnPresence(func(e centrifuge.PresenceEvent, callback centrifuge.PresenceCallback) {
			callback(PresenceHandler(node, client, &e))
		})
		client.OnPresenceStats(func(e centrifuge.PresenceStatsEvent, callback centrifuge.PresenceStatsCallback) {
			callback(PresenceStatsHandler(node, client, &e))
		})
		client.OnHistory(func(e centrifuge.HistoryEvent, callback centrifuge.HistoryCallback) {
			callback(HistoryHandler(node, client, &e))
		})
	})
}
