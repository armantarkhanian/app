// Package jwt ..
package jwt

import (
	"app/internal/pkg/redis"
	"errors"
	"log"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	AuthMiddleware *jwt.GinJWTMiddleware
)

func LoginHandler() gin.HandlerFunc {
	return AuthMiddleware.LoginHandler
}

func GetPayload(c *gin.Context) (userID string, sessionID string) {
	claims := jwt.ExtractClaims(c)
	userID, _ = claims["sub"].(string)
	sessionID, _ = claims["sessionID"].(string)
	return userID, sessionID
}

type payload struct {
	UserID    string `json:"userID"`
	SessionID string `json:"sessionID"`
}

func init() {
	var err error
	AuthMiddleware, err = jwt.New(&jwt.GinJWTMiddleware{
		Key:     []byte("945192e8f34a44b2af60265df7860e79898ec86108cc4b52a83bbc19498001ad"),
		Timeout: time.Hour * 24 * 90,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*payload); ok {
				return jwt.MapClaims{
					"sub":       v.UserID,
					"sessionID": v.SessionID,
				}
			}
			return jwt.MapClaims{}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			userID := "arman"
			sessionID := uuid.New().String()

			if err := redis.Client.SAdd(redis.Ctx, "u:"+userID+":s", sessionID).Err(); err != nil {
				log.Println("[ERROR] [redis]", err)
				return nil, errors.New("database error")
			}

			return &payload{
				UserID:    userID,
				SessionID: sessionID,
			}, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(*payload); ok {
				cmd := redis.Client.SIsMember(redis.Ctx, "u:"+v.UserID+":s", v.SessionID)
				if cmd.Err() != nil {
					log.Println("[ERROR] [redis]", cmd.Err())
				}
				return cmd.Val()
			}
			return false
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			var data payload
			data.UserID, _ = claims["sub"].(string)
			data.SessionID, _ = claims["sessionID"].(string)
			return &data
		},
		LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			c.JSON(code, gin.H{
				"statusCode": code,
				"token":      token,
				"expire":     expire.Unix(),
			})
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"error": "unauthorized",
			})
		},
		TokenLookup:   "header: Authorization, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now().UTC,

		SendCookie:     true,
		SecureCookie:   false,
		CookieHTTPOnly: true,
		CookieDomain:   "localhost",
		CookieName:     "jwt",
		CookieSameSite: http.SameSiteStrictMode,
	})

	if err != nil {
		panic("JWT Error:" + err.Error())
	}
}
