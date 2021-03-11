// Package middlewares ...
package middlewares

import (
	"app/internal/pkg/geoip"
	"app/internal/pkg/global"
	"app/internal/pkg/keys"
	internalSessions "app/internal/pkg/sessions"
	"fmt"
	"log"
	"strings"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Sessions() gin.HandlerFunc {
	return internalSessions.Middleware()
}

func GeoIP() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(keys.CountryCode, geoip.CountryCodeByIP(c.ClientIP()))
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
		claims := jwt.ExtractClaims(c)
		userID, _ := claims["userID"].(string)

		start := time.Now().UTC()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		timeStamp := time.Now().UTC()
		latency := timeStamp.Sub(start)

		if raw != "" {
			path = path + "?" + raw
		}

		v, _ := c.Get(keys.CountryCode)
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

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		sessionID, _ := session.Get(keys.SessionID).(string)
		userID, _ := session.Get(keys.UserID).(string)
		if userID == "" || sessionID == "" {
			c.JSON(200, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}
