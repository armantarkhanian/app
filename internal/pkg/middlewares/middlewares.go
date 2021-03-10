// Package middlewares ...
package middlewares

import (
	"app/internal/pkg/configs"
	"app/internal/pkg/geoip"
	"app/internal/pkg/global"
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
			if r := recover(); r != nil {
				array := strings.Split(c.HandlerName(), "/")
				handler := fmt.Sprintf("in Handler: %q", array[len(array)-1])
				log.Println("[PANIC]", r, handler)
				c.HTML(200, "panic.html", nil)
				c.Abort()
			}
		}()
		c.Next()
	}
}

func AccessLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID, _ := session.Get("userID").(string)
		start := time.Now().UTC()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		timeStamp := time.Now().UTC()
		latency := timeStamp.Sub(start)

		if raw != "" {
			path = path + "?" + raw
		}

		v, _ := c.Get("countryCode")
		countryCode := v.(string)

		errs := c.Errors.ByType(gin.ErrorTypePrivate).String()
		if errs != "" {
			log.Println("[ERROR]", errs)
		}

		global.AccessLog.Log().
			Int64("time", timeStamp.Unix()).
			Int("backendID", global.BackendID).
			Str("countryCode", countryCode).
			Str("ip", c.ClientIP()).
			Str("userID", userID).
			Str("userAgent", c.Request.UserAgent()).
			Int64("latency", latency.Milliseconds()).
			Int("statusCode", c.Writer.Status()).
			Str("method", c.Request.Method).
			Str("path", path).
			Int64("requestBodySize", c.Request.ContentLength).
			Int("responseBodySize", c.Writer.Size()).
			Str("errors", errs).
			Msg("")
	}
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

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID, _ := session.Get("userID").(string)
		if userID == "" {
			c.JSON(200, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}
