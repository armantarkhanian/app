// Package configs ...
package configs

type Configs struct {
	Sessions sessionsConfig `json:"sessions"`
	Gin      ginConfig      `json:"gin"`
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
	Addr     string   `json:"addr" validate:"required"`
	Timeouts Timeouts `json:"timeouts" validate:"required"`
	Mode     string   `json:"mode" validate:"required,oneof=test debug release"`
}

type Timeouts struct {
	Read     string `json:"read" validate:"required"`
	Write    string `json:"write" validate:"required"`
	Idle     string `json:"idle" validate:"required"`
	Shutdown string `json:"shutdown" validate:"required"`
}
