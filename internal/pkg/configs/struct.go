// Package configs ...
package configs

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type Configs struct {
	Notify   notify         `json:"notify"`
	Sessions sessionsConfig `json:"sessions"`
	Gin      ginConfig      `json:"gin"`
	Redis    redisConfig    `json:"redis"`
}

type redisConfig struct {
	Mode string   `json:"mode"`
	Addr []string `json:"addr"`
}

type notify struct {
	Telegram telegramNotify `json:"telegram"`
}

type telegramNotify struct {
	BotToken string `json:"botToken"`
	ChatID   string `json:"chatID"`
}

func (telegramWriter *telegramNotify) Write(p []byte) (n int, err error) {
	str := string(p)
	fmt.Println(telegramWriter)
	baseURL, err := url.Parse("https://api.telegram.org/bot" + telegramWriter.BotToken + "/sendMessage")
	if err != nil {
		log.Println("[telegram]", err)
		return 0, err
	}

	params := url.Values{}

	params.Add("chat_id", telegramWriter.ChatID)
	params.Add("text", str)

	baseURL.RawQuery = params.Encode()

	link := baseURL.String()

	resp, err := http.Get(link)
	if err != nil {
		log.Println("[telegram]", err)
		return 0, err
	}
	defer resp.Body.Close()

	if err != nil {
		log.Println("[telegram]", err)
		return 0, err
	}

	return 0, resp.Body.Close()
}

type sessionsConfig struct {
	Domain string `json:"domain" validate:"required"`
	MaxAge int    `json:"maxAge" validate:"required"`
	Secure bool   `json:"secure"`
}

type ginConfig struct {
	Addr                       string      `json:"addr" validate:"required"`
	Timeouts                   timeouts    `json:"timeouts" validate:"required"`
	Middlewares                middlewares `json:"middlewares"`
	Mode                       string      `json:"mode" validate:"required,oneof=test debug release"`
	ConsoleLog                 bool        `json:"consoleLog"`
	AccessLoggerTimeLayout     string      `json:"accessLoggerTimeLayout"`
	QueriesPerMinuteForCaptcha int         `json:"queriesPerMinuteForCaptcha"`
}

type middlewares struct {
	Recovery     bool `json:"recovery"`
	AccessLogger bool `json:"accessLogger"`
	JWT          bool `json:"jwt"`
	GeoIP        bool `json:"geoip"`
	Sessions     bool `json:"sessions"`
}

type timeouts struct {
	Read     string `json:"read" validate:"required"`
	Write    string `json:"write" validate:"required"`
	Idle     string `json:"idle" validate:"required"`
	Shutdown string `json:"shutdown" validate:"required"`
}
