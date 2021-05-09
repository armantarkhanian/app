// Package websocket ...
package websocket

import (
	"app/internal/pkg/logger"
	"fmt"
	"log"

	"github.com/centrifugal/centrifuge"
)

func AliveHandler(c *centrifuge.Client) {}
func DisconectHandler(c *centrifuge.Client, e *centrifuge.DisconnectEvent) {
	log.Printf("client %q disconnected", c.UserID())
}

func SubscribeHandler(c *centrifuge.Client, e *centrifuge.SubscribeEvent) (centrifuge.SubscribeReply, error) {
	return centrifuge.SubscribeReply{}, nil
}
func UnsubscribeHandler(c *centrifuge.Client, e *centrifuge.UnsubscribeEvent) {}

func PublishHandler(c *centrifuge.Client, e *centrifuge.PublishEvent) (centrifuge.PublishReply, error) {
	// запретить пользователям публиковать данные
	return centrifuge.PublishReply{}, centrifuge.ErrorPermissionDenied
}

func RefreshHandler(c *centrifuge.Client, e *centrifuge.RefreshEvent) (centrifuge.RefreshReply, error) {
	return centrifuge.RefreshReply{}, nil
}

func SubRefreshHandler(c *centrifuge.Client, e *centrifuge.SubRefreshEvent) (centrifuge.SubRefreshReply, error) {
	return centrifuge.SubRefreshReply{}, nil
}

func RPCHandler(c *centrifuge.Client, e *centrifuge.RPCEvent) (centrifuge.RPCReply, error) {
	switch e.Method {
	case "stopTyping":
		data := `
{
"action": "stopTyping",
"user": "` + c.UserID() + `"
}
`
		if _, err := node.Publish("chat", []byte(data)); err != nil {
			logger.Error(err)
		}
	case "typing":
		data := `
{
"action": "typing",
"user": "` + c.UserID() + `"
}
`
		if _, err := node.Publish("chat", []byte(data)); err != nil {
			logger.Error(err)
		}
	case "sendMessage":
		if string(e.Data) == "logout" {
			if _, err := node.Publish("#"+c.UserID(), []byte("1")); err != nil {
				logger.Error(err)
			}
		}
		fmt.Println(e.Method, string(e.Data), c.UserID())
		node.Publish("chat", []byte(`{"user": "`+c.UserID()+`", "action":"sendMessage", "data":"`+string(e.Data)+`"}`))
	default:
		return centrifuge.RPCReply{}, centrifuge.ErrorMethodNotFound
	}
	return centrifuge.RPCReply{}, nil
}

func MessageHandler(c *centrifuge.Client, e *centrifuge.MessageEvent) {
	c.Send(e.Data)
}

func PresenceHandler(c *centrifuge.Client, e *centrifuge.PresenceEvent) (centrifuge.PresenceReply, error) {
	return centrifuge.PresenceReply{}, nil
}

func PresenceStatsHandler(c *centrifuge.Client, e *centrifuge.PresenceStatsEvent) (centrifuge.PresenceStatsReply, error) {
	return centrifuge.PresenceStatsReply{}, nil
}

func HistoryHandler(c *centrifuge.Client, e *centrifuge.HistoryEvent) (centrifuge.HistoryReply, error) {
	return centrifuge.HistoryReply{}, nil
}

func IsUserOnline(userID string) (bool, error) {
	return checkUserOnline(userID)
}

func checkUserOnline(userID string) (bool, error) {
	res1, _ := node.Presence("#" + userID)
	for key, value := range res1.Presence {
		fmt.Println(key, "=>", value)
	}
	res, err := node.PresenceStats("#" + userID)
	if err != nil {
		return false, err
	}
	return res.NumClients > 0, nil
}
