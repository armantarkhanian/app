// Package jwt ..
package jwt

import (
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

var (
	AuthMiddleware *jwt.GinJWTMiddleware
)

type payload struct {
	UserID    string
	SessionID string
}

func init() {
	var err error
	AuthMiddleware, err = jwt.New(&jwt.GinJWTMiddleware{
		Realm:   "test zone",
		Key:     []byte("4aceccc4ae3d43e28e7788c6165105e0"),
		Timeout: time.Hour * 24 * 90,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*payload); ok {
				return jwt.MapClaims{
					"userID":    v.UserID,
					"sessionID": v.SessionID,
				}
			}
			return jwt.MapClaims{}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			return &payload{
				UserID:    "arman",
				SessionID: "1234",
			}, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(*payload); ok {
				// check is v.SessionID is invalid
				return v.SessionID == "123"
			}
			return false
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &payload{
				UserID:    claims["userID"].(string),
				SessionID: claims["sessionID"].(string),
			}
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

	errInit := AuthMiddleware.MiddlewareInit()

	if errInit != nil {
		panic("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}

}
