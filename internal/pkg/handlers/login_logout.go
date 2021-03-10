package handlers

import (
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func login(c *gin.Context) {
	session := sessions.Default(c)
	session.Set("userID", "admin")
	session.Set("cookieID", 15)
	if err := session.Save(); err != nil {
		log.Println("[ERROR]", err)
		c.JSON(200, gin.H{"ok": false})
	}
	c.JSON(200, gin.H{"ok": true})
}

func logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Set("userID", "")
	if err := session.Save(); err != nil {
		log.Println("[ERROR]", err)
		c.JSON(200, gin.H{"ok": false})
	}
	c.JSON(200, gin.H{"ok": true})
}
