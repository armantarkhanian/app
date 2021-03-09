package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createAccountRequestData1 struct {
	Username     string `json:"username"`
	AccountName  string `json:"accountName"`
	Password     string `json:"password"`
	Email        string `json:"email"`
	IsUser       bool   `json:"isUser"`
	DayOfBirth   string `json:"dayOfBirth"`
	MonthOfBirth string `json:"monthOfBirth"`
	YearOfBirth  string `json:"yearOfBirth"`
}

func createAccount1(c *gin.Context) {
	jsonData, _ := c.MustGet("requestData").(string)
	var requestData createAccountRequestData1
	if err := json.Unmarshal([]byte(jsonData), &requestData); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, requestData)
}
