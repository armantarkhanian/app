// Package middlewares ...
package middlewares

import (
	"app/internal/pkg/geoip"
	"app/internal/pkg/global"
	internalJWT "app/internal/pkg/jwt"
	"app/internal/pkg/keys"
	"app/internal/pkg/metrics"
	internalSessions "app/internal/pkg/sessions"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	ua "github.com/mileusna/useragent"
)

func Sessions() gin.HandlerFunc {
	return internalSessions.Middleware()
}

func Auth() gin.HandlerFunc {
	return internalJWT.AuthMiddleware.MiddlewareFunc()
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
		userID, _ := internalJWT.GetPayload(c)

		start := time.Now().UTC()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		timeStamp := time.Now().UTC()
		latency := timeStamp.Sub(start)

		if raw != "" {
			path = path + "?" + raw
		}

		var countryCode string
		v, _ := c.Get(keys.CountryCode)
		if v != nil {
			countryCode, _ = v.(string)
		}

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

func Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		path := c.Request.URL.Path

		ua := ua.Parse(c.Request.UserAgent())
		if ua.OS == "" {
			ua.OS = "others"
		}
		metrics.Incr("handlers", path)
		metrics.Incr("methods", c.Request.Method)
		metrics.Incr("statusCodes", strconv.Itoa(c.Writer.Status()))
		metrics.Incr("os", ua.OS)
		metrics.Incr("browsers", ua.Name)

		var ip string
		if c.ClientIP() == "::1" {
			ip = "79.141.162.81"
		} else {
			ip = "51.83.70.23"
		}
		countryCode := geoip.CountryCodeByIP(ip)
		if countryCode == "" {
			countryCode = "others"
		}
		metrics.Incr("countries", countryCode)
	}
}
