// Package recaptcha ...
package recaptcha

import (
	"app/internal/pkg/configs"
	"app/internal/pkg/httpclient"
	"encoding/json"
	"log"
	"time"

	"github.com/gin-contrib/sessions"

	"github.com/gin-gonic/gin"
)

var (
	secret = "6Le_BHgaAAAAAPpa6p90cH16r3Wn51G21HqAMd5v"
)

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if configs.Store.Gin.QueriesPerMinuteForCaptcha <= 0 {
			c.Next()
			return
		}
		session := sessions.Default(c)
		needCaptcha, _ := session.Get("needCaptcha").(bool)
		if needCaptcha {
			c.IndentedJSON(200, gin.H{
				"error":       "too many queries",
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
					"error":       "too many queries",
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

// CheckCaptcha gets POST request with "g-recaptcha-response" and validate it
func CheckCaptcha(c *gin.Context) {
	session := sessions.Default(c)
	needCaptch, _ := session.Get("needCaptcha").(bool)
	if !needCaptch {
		c.JSON(200, gin.H{
			"ok": true,
		})
		return
	}
	type response struct {
		Success    bool     `json:"success"`
		Hostname   string   `json:"hostname"`
		ErrorCodes []string `json:"error-codes"`
	}

	resp, err := httpclient.Post("https://www.google.com/recaptcha/api/siteverify", map[string]string{
		"response": c.PostForm("g-recaptcha-response"),
		"secret":   secret,
		"remoteip": c.ClientIP(),
	})
	if err != nil {
		log.Println("[ERROR]", err)
		c.JSON(200, gin.H{
			"ok": false,
		})
		return
	}
	var respStruct response
	err = json.Unmarshal(resp, &respStruct)
	if err != nil {
		log.Println("[ERROR]", err)
		c.JSON(200, gin.H{
			"ok": false,
		})
		return
	}
	if respStruct.Success {
		session.Set("lastActionTime", time.Now().UTC())
		session.Set("qps", 0)
		session.Set("needCaptcha", false)
		if err := session.Save(); err != nil {
			log.Println("[ERROR]", err)
			c.JSON(200, gin.H{
				"ok": false,
			})
		}
		c.JSON(200, gin.H{
			"ok": true,
		})
		return
	}
	c.JSON(200, gin.H{
		"ok": false,
	})
}
