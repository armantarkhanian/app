// Package middlewares ...
package middlewares

import (
	"app/internal/pkg/configs"
	"app/internal/pkg/geoip"
	internalSessions "app/internal/pkg/sessions"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Sessions() gin.HandlerFunc {
	return internalSessions.Middleware()
}

func GeoIP() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("countryCode", geoip.CountryCodeByIP(c.ClientIP()))
		c.Next()
	}
}

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				array := strings.Split(c.HandlerName(), "/")
				log.Println("[PANIC]", array[len(array)-1], err)
				c.HTML(200, "panic.html", nil)
				c.Abort()
			}
		}()
		c.Next()
	}
}

func AccessLogger() gin.HandlerFunc {
	return gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: func(param gin.LogFormatterParams) string {
			if param.Latency > time.Minute {
				param.Latency = param.Latency - param.Latency%time.Second
			}

			if param.ErrorMessage != "" {
				param.ErrorMessage = fmt.Sprintf(`| %s",`, param.ErrorMessage)
			}

			return fmt.Sprintf(`%v | %15s | %13v | %3d | %-9s | %s %s`,
				param.TimeStamp.Format("02 Jan 2006 15:04:05 MST"),
				param.ClientIP,
				param.Latency,
				param.StatusCode,
				"["+param.Method+"]",
				param.Path,
				param.ErrorMessage,
			) + "\n"
		},
	})
}

func RecaptchaProtect() gin.HandlerFunc {
	return func(c *gin.Context) {
		if configs.Store.Gin.QueriesPerMinuteForCaptcha <= 0 {
			c.Next()
			return
		}
		session := sessions.Default(c)
		needCaptcha, _ := session.Get("needCaptcha").(bool)
		if needCaptcha {
			c.IndentedJSON(200, gin.H{
				"needCaptcha": true,
			})
			c.Abort()
			return
		}
		lastActionTime, _ := session.Get("lastActionTime").(time.Time)
		qps, _ := session.Get("qps").(int)

		if lastActionTime.IsZero() {
			session.Set("lastActionTime", time.Now().UTC())
			if err := session.Save(); err != nil {
				log.Println("[ERROR]", err)
			}
			c.Next()
			return
		}

		ts := time.Now().UTC().Sub(lastActionTime).Milliseconds()
		if ts <= 60000 {
			qps++
			if qps >= configs.Store.Gin.QueriesPerMinuteForCaptcha {
				session.Set("needCaptcha", true)
				if err := session.Save(); err != nil {
					log.Println("[ERROR]", err)
				}
				c.IndentedJSON(200, gin.H{
					"needCaptcha": true,
				})
				c.Abort()
			} else {
				session.Set("qps", qps)
				if err := session.Save(); err != nil {
					log.Println("[ERROR]", err)
				}
			}
			return
		}

		session.Set("lastActionTime", time.Now().UTC())
		session.Set("qps", 0)
		if err := session.Save(); err != nil {
			log.Println("[ERROR]", err)
		}
		c.Next()
	}
}
