// Package configs ...
package configs

type Configs struct {
	Notify   notify         `json:"notify"`
	Sessions sessionsConfig `json:"sessions"`
	Gin      ginConfig      `json:"gin"`
}

type notify struct {
	Telegram telegramNotify `json:"telegram"`
}

type telegramNotify struct {
	BotToken string `json:"botToken"`
	ChatID   string `json:"chatID"`
}

type sessionsConfig struct {
	Domain string `json:"domain" validate:"required"`
	MaxAge int    `json:"maxAge" validate:"required"`
	Secure bool   `json:"secure"`
}

type ginConfig struct {
	Addr     string   `json:"addr" validate:"required"`
	Timeouts timeouts `json:"timeouts" validate:"required"`
	Mode     string   `json:"mode" validate:"required,oneof=test debug release"`
}

type timeouts struct {
	Read     string `json:"read" validate:"required"`
	Write    string `json:"write" validate:"required"`
	Idle     string `json:"idle" validate:"required"`
	Shutdown string `json:"shutdown" validate:"required"`
}
