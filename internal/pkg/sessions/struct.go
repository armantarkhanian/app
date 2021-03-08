// Package sessions ...
package sessions

import "time"

type Session struct {
	UserID           string
	Username         string
	LastActionTime   time.Time
	ActionsPerSecond int
}
