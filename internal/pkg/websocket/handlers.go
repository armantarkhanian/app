// Package websocket ...
package websocket

import (
	"app/internal/pkg/logger"
	"context"
	"fmt"

	"github.com/centrifugal/centrifuge"
)

func setHandlers(node *centrifuge.Node) {
	node.OnNotification(func(e centrifuge.NotificationEvent) {
		fmt.Println(e.Op, string(e.Data), e.FromNodeID)
	})
	node.OnConnecting(func(ctx context.Context, e centrifuge.ConnectEvent) (centrifuge.ConnectReply, error) {
		cred, ok := centrifuge.GetCredentials(ctx)
		if !ok {
			return centrifuge.ConnectReply{}, centrifuge.ErrorBadRequest
		}
		if cred.UserID == "" {
			return centrifuge.ConnectReply{}, centrifuge.ErrorBadRequest
		}
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
		logger.Infof("%q connected via %q using %q protocl", client.UserID(), transportName, transportProto)
		client.OnAlive(func() {
			AliveHandler(client)
		})
		client.OnDisconnect(func(e centrifuge.DisconnectEvent) {
			DisconectHandler(client, &e)
		})
		client.OnSubscribe(func(e centrifuge.SubscribeEvent, callback centrifuge.SubscribeCallback) {
			callback(SubscribeHandler(client, &e))
		})
		client.OnUnsubscribe(func(e centrifuge.UnsubscribeEvent) {
			UnsubscribeHandler(client, &e)
		})

		/* not allow users to publish directly to channel
		client.OnPublish(func(e centrifuge.PublishEvent, callback centrifuge.PublishCallback) {
			callback(PublishHandler(client, &e))
		})*/

		client.OnRefresh(func(e centrifuge.RefreshEvent, callback centrifuge.RefreshCallback) {
			callback(RefreshHandler(client, &e))
		})
		client.OnSubRefresh(func(e centrifuge.SubRefreshEvent, callback centrifuge.SubRefreshCallback) {
			callback(SubRefreshHandler(client, &e))
		})
		client.OnRPC(func(e centrifuge.RPCEvent, callback centrifuge.RPCCallback) {
			callback(RPCHandler(client, &e))
		})
		client.OnMessage(func(e centrifuge.MessageEvent) {
			MessageHandler(client, &e)
		})
		client.OnPresence(func(e centrifuge.PresenceEvent, callback centrifuge.PresenceCallback) {
			callback(PresenceHandler(client, &e))
		})
		client.OnPresenceStats(func(e centrifuge.PresenceStatsEvent, callback centrifuge.PresenceStatsCallback) {
			callback(PresenceStatsHandler(client, &e))
		})
		client.OnHistory(func(e centrifuge.HistoryEvent, callback centrifuge.HistoryCallback) {
			callback(HistoryHandler(client, &e))
		})
	})
}
