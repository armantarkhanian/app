// Package recaptcha ...
package recaptcha

import (
	"app/internal/pkg/configs"
	"log"
	"time"

	"github.com/gin-contrib/sessions"

	"github.com/gin-gonic/gin"
)

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		needCaptcha, _ := session.Get("needCaptcha").(bool)
		if needCaptcha {
			c.IndentedJSON(200, gin.H{
				"error": "too many queries per second",
			})
			c.Abort()
			return
		}
		lastActionTime, _ := session.Get("lastActionTime").(time.Time)
		qps, _ := session.Get("qps").(int)

		if configs.Store.Gin.QueriesPerMinuteForCaptcha <= 0 {
			c.Next()
			return
		}

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
					"error": "too many queries per second",
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
