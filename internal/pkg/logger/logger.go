// Package logger ...
package logger

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("[PANIC]", c.HandlerName(), err)
				c.HTML(200, "panic.html", nil)
				c.Abort()
			}
		}()
		c.Next()
	}
}

func Middleware() gin.HandlerFunc {
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
				param.Request.URL.Path,
				param.ErrorMessage,
			) + "\n"
		},
	})
}
