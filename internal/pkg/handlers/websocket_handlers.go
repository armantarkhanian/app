package handlers

import (
	"fmt"
	"log"

	"github.com/centrifugal/centrifuge"
)

type WebSocket struct {
	aliveHandler         centrifuge.AliveHandler
	disconnectHandler    centrifuge.DisconnectHandler
	subscribeHandler     centrifuge.SubscribeHandler
	unsubscribeHandler   centrifuge.UnsubscribeHandler
	publishHandler       centrifuge.PublishHandler
	refreshHandler       centrifuge.RefreshHandler
	subRefreshHandler    centrifuge.SubRefreshHandler
	rpcHandler           centrifuge.RPCHandler
	messageHandler       centrifuge.MessageHandler
	presenceHandler      centrifuge.PresenceHandler
	presenceStatsHandler centrifuge.PresenceStatsHandler
	historyHandler       centrifuge.HistoryHandler
}

func RPCHandler(n *centrifuge.Node, c *centrifuge.Client, e *centrifuge.RPCEvent) (centrifuge.RPCReply, error) {
	return centrifuge.RPCReply{Data: e.Data}, nil
}

func onSubscribe(n *centrifuge.Node, c *centrifuge.Client, e *centrifuge.SubscribeEvent) (centrifuge.SubscribeReply, error) {
	return centrifuge.SubscribeReply{}, nil
}

func onMessage(n *centrifuge.Node, c *centrifuge.Client, e *centrifuge.MessageEvent) {
	c.Send(e.Data)
}

func onDisconect(n *centrifuge.Node, c *centrifuge.Client, e *centrifuge.DisconnectEvent) {
	log.Printf("client %q disconnected", c.UserID())
}

func onPublish(n *centrifuge.Node, c *centrifuge.Client, e *centrifuge.PublishEvent) (centrifuge.PublishReply, error) {
	if string(e.Data) == `"logout"` {
		fmt.Println("Do it")
		fmt.Println(n.Publish("user_15", []byte("1")))
	}
	if c.UserID() != "user_15" {
		c.Disconnect(&centrifuge.Disconnect{
			Code:      200,
			Reason:    "you do not have *permissions here",
			Reconnect: false,
		})
	}

	return centrifuge.PublishReply{}, nil
}
