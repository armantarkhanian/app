package handlers

import (
	"fmt"
	"log"

	"github.com/centrifugal/centrifuge"
)

func PresenceStatsHandler(n *centrifuge.Node, c *centrifuge.Client, e *centrifuge.PresenceStatsEvent) (centrifuge.PresenceStatsReply, error) {
	return centrifuge.PresenceStatsReply{}, nil
}
func RefreshHandler(n *centrifuge.Node, c *centrifuge.Client, e *centrifuge.RefreshEvent) (centrifuge.RefreshReply, error) {
	return centrifuge.RefreshReply{}, nil
}
func SubRefreshHandler(n *centrifuge.Node, c *centrifuge.Client, e *centrifuge.SubRefreshEvent) (centrifuge.SubRefreshReply, error) {
	return centrifuge.SubRefreshReply{}, nil
}

func PresenceHandler(n *centrifuge.Node, c *centrifuge.Client, e *centrifuge.PresenceEvent) (centrifuge.PresenceReply, error) {
	return centrifuge.PresenceReply{}, nil
}
func HistoryHandler(n *centrifuge.Node, c *centrifuge.Client, e *centrifuge.HistoryEvent) (centrifuge.HistoryReply, error) {
	return centrifuge.HistoryReply{}, nil
}

func AliveHandler(n *centrifuge.Node, c *centrifuge.Client)                                       {}
func UnsubscribeHandler(n *centrifuge.Node, c *centrifuge.Client, e *centrifuge.UnsubscribeEvent) {}

func RPCHandler(n *centrifuge.Node, c *centrifuge.Client, e *centrifuge.RPCEvent) (centrifuge.RPCReply, error) {
	return centrifuge.RPCReply{}, nil
}

func SubscribeHandler(n *centrifuge.Node, c *centrifuge.Client, e *centrifuge.SubscribeEvent) (centrifuge.SubscribeReply, error) {
	return centrifuge.SubscribeReply{}, nil
}

func MessageHandler(n *centrifuge.Node, c *centrifuge.Client, e *centrifuge.MessageEvent) {
	c.Send(e.Data)
}

func DisconectHandler(n *centrifuge.Node, c *centrifuge.Client, e *centrifuge.DisconnectEvent) {
	log.Printf("client %q disconnected", c.UserID())
}

func PublishHandler(n *centrifuge.Node, c *centrifuge.Client, e *centrifuge.PublishEvent) (centrifuge.PublishReply, error) {
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
