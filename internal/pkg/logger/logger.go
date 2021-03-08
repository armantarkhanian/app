// Package logger ...
package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func Init() {
	logFile, err := os.OpenFile("./logs/main.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("[FATAL] error opening file:", err)
	}
	accessLog, err := os.OpenFile("./logs/http.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("[FATAL] error opening file:", err)
	}
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.LUTC)
	log.SetOutput(logFile)

	gin.DefaultWriter = accessLog
	gin.DefaultErrorWriter = logFile
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
				param.Path,
				param.ErrorMessage,
			) + "\n"
		},
	})
}
