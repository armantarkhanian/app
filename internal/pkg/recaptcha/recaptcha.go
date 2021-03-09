// Package recaptcha ...
package recaptcha

import (
	"app/internal/pkg/httpclient"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"

	"github.com/gin-gonic/gin"
)

var (
	secret = "6Le_BHgaAAAAAPpa6p90cH16r3Wn51G21HqAMd5v"
)

type verifyCaptchaData1 struct {
	RecaptchaResponse string `json:"recaptchaResponse"`
}

type googleResponse struct {
	Success    bool     `json:"success"`
	Hostname   string   `json:"hostname"`
	ErrorCodes []string `json:"error-codes"`
}

func VerifyCaptcha(c *gin.Context) {
	session := sessions.Default(c)
	needCaptch, _ := session.Get("needCaptcha").(bool)
	if !needCaptch {
		c.JSON(200, gin.H{
			"ok": true,
		})
		return
	}
	jsonData, _ := c.MustGet("requestData").(string)
	var requestData verifyCaptchaData1
	if err := json.Unmarshal([]byte(jsonData), &requestData); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	resp, err := httpclient.Post("https://www.google.com/recaptcha/api/siteverify", map[string]string{
		"response": requestData.RecaptchaResponse,
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

	var respStruct googleResponse

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
