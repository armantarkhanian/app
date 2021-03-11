// Package jwt ..
package jwt

import (
	"app/internal/pkg/keys"
	"fmt"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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
	Hash    string `json:"hash"`
}

func hash(str string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.MinCost)
	return string(bytes), err
}

func checkHash(str, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(str))
	return err == nil
}

func init() {
	var err error
	AuthMiddleware, err = jwt.New(&jwt.GinJWTMiddleware{
		Realm:            "test zone",
		SigningAlgorithm: "HS256",
		Key:              []byte("4aceccc4ae3d43e28e7788c6165105e0"),
		Timeout:          time.Hour * 24 * 90,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*payload); ok {
				return jwt.MapClaims{
					keys.UserID:  v.UserID,
					keys.TokenID: v.TokenID,
					"hash":       v.Hash,
				}
			}
			return jwt.MapClaims{}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			hs, err := hash(c.ClientIP() + c.Request.UserAgent())
			if err != nil {
				panic(err)
			}
			return &payload{
				UserID:  "arman",
				TokenID: uuid.New().String(),
				Hash:    hs,
			}, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(*payload); ok {
				if !checkHash(c.ClientIP()+c.Request.UserAgent(), v.Hash) {
					fmt.Println("Invalid hash") // send message to user about it
				}
				return true
			}
			return false
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &payload{
				UserID:  claims[keys.UserID].(string),
				TokenID: claims[keys.TokenID].(string),
				Hash:    claims["hash"].(string),
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
