// Package jwt ..
package jwt

import (
	"app/internal/pkg/keys"
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

func GetPayload(c *gin.Context) (userID string, tokenID string) {
	claims := jwt.ExtractClaims(c)
	userID, _ = claims[keys.UserID].(string)
	tokenID, _ = claims[keys.TokenID].(string)
	return userID, tokenID
}

type payload struct {
	UserID  string `json:"userID"`
	TokenID string `json:"tokenID"`
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
					keys.UserID:  v.UserID,
					keys.TokenID: v.TokenID,
				}
			}
			return jwt.MapClaims{}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			return &payload{
				UserID:  "arman",
				TokenID: uuid.New().String(),
			}, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			// here check if tokenID is invalid
			return true
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &payload{
				UserID:  claims[keys.UserID].(string),
				TokenID: claims[keys.TokenID].(string),
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
