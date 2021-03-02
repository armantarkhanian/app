// Package configs ...
package configs

type Configs struct {
	TelegramNotify telegramNotify `json:"telegramNotify"`
	Sessions       sessionsConfig `json:"sessions"`
	Gin            ginConfig      `json:"gin"`
}

type telegramNotify struct {
	BotToken string `json:"botToken"`
	ChatID   string `json:"chatID"`
}

type sessionsConfig struct {
	Name              string `json:"name" validate:"required,alphanum"`
	AuthenticationKey string `json:"authenticationKey" validate:"required,alphanum,len=32"`
	EncryptionKey     string `json:"encryptionKey" validate:"required,alphanum,len=32"`
	Domain            string `json:"domain" validate:"required"`
	MaxAge            int    `json:"maxAge" validate:"required"`
	Secure            bool   `json:"secure"`
}

type ginConfig struct {
	Addr     string    `json:"addr" validate:"required"`
	Timeouts timeouts  `json:"timeouts" validate:"required"`
	Mode     string    `json:"mode" validate:"required,oneof=test debug release"`
	Log      logConfig `json:"log"`
}

type logConfig struct {
	AccessLogFile string `json:"accessLogFile"`
	ErrorLogFile  string `json:"errorLogFile"`
	UseStdOut     bool   `json:"useStdOut"`
	UseStdErr     bool   `json:"userStdErr"`
}

type timeouts struct {
	Read     string `json:"read" validate:"required"`
	Write    string `json:"write" validate:"required"`
	Idle     string `json:"idle" validate:"required"`
	Shutdown string `json:"shutdown" validate:"required"`
}
