// Package websocket ...
package websocket

import (
	"app/internal/pkg/logger"
	"fmt"
	"log"

	"github.com/centrifugal/centrifuge"
)

func AliveHandler(n *centrifuge.Node, c *centrifuge.Client) {}
func DisconectHandler(n *centrifuge.Node, c *centrifuge.Client, e *centrifuge.DisconnectEvent) {
	log.Printf("client %q disconnected", c.UserID())
}

func SubscribeHandler(n *centrifuge.Node, c *centrifuge.Client, e *centrifuge.SubscribeEvent) (centrifuge.SubscribeReply, error) {
	return centrifuge.SubscribeReply{}, nil
}
func UnsubscribeHandler(n *centrifuge.Node, c *centrifuge.Client, e *centrifuge.UnsubscribeEvent) {}

func PublishHandler(n *centrifuge.Node, c *centrifuge.Client, e *centrifuge.PublishEvent) (centrifuge.PublishReply, error) {
	fmt.Println(string(e.Data))
	if string(e.Data) == "logout" {
		if _, err := n.Publish("#user_15", []byte("1")); err != nil {
			logger.Error(err)
		}
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

func RefreshHandler(n *centrifuge.Node, c *centrifuge.Client, e *centrifuge.RefreshEvent) (centrifuge.RefreshReply, error) {
	return centrifuge.RefreshReply{}, nil
}

func SubRefreshHandler(n *centrifuge.Node, c *centrifuge.Client, e *centrifuge.SubRefreshEvent) (centrifuge.SubRefreshReply, error) {
	return centrifuge.SubRefreshReply{}, nil
}

func RPCHandler(n *centrifuge.Node, c *centrifuge.Client, e *centrifuge.RPCEvent) (centrifuge.RPCReply, error) {
	return centrifuge.RPCReply{}, nil
}

func MessageHandler(n *centrifuge.Node, c *centrifuge.Client, e *centrifuge.MessageEvent) {
	c.Send(e.Data)
}

func PresenceHandler(n *centrifuge.Node, c *centrifuge.Client, e *centrifuge.PresenceEvent) (centrifuge.PresenceReply, error) {
	return centrifuge.PresenceReply{}, nil
}

func PresenceStatsHandler(n *centrifuge.Node, c *centrifuge.Client, e *centrifuge.PresenceStatsEvent) (centrifuge.PresenceStatsReply, error) {
	return centrifuge.PresenceStatsReply{}, nil
}

func HistoryHandler(n *centrifuge.Node, c *centrifuge.Client, e *centrifuge.HistoryEvent) (centrifuge.HistoryReply, error) {
	return centrifuge.HistoryReply{}, nil
}
