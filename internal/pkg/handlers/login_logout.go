package handlers

import (
	"app/internal/pkg/keys"
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func login(c *gin.Context) {
	session := sessions.Default(c)
	session.Set(keys.SessionID, "sessionID 3324")
	session.Set(keys.UserID, "Admin")
	session.Set(keys.Username, "Fuck you bitch")
	if err := session.Save(); err != nil {
		log.Println("[ERROR]", err)
		c.JSON(200, gin.H{"ok": false})
	}
	c.JSON(200, gin.H{"ok": true})
}

func logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Set(keys.UserID, "")
	if err := session.Save(); err != nil {
		log.Println("[ERROR]", err)
		c.JSON(200, gin.H{"ok": false})
	}
	c.JSON(200, gin.H{"ok": true})
}
