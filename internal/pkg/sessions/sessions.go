/* Package sessions is ... */
package sessions

import (
	"app/internal/pkg/configs"
	"encoding/gob"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

var store cookie.Store

var sessionName string

func Init() {
	gob.Register(Session{})

	store = cookie.NewStore(
		[]byte(configs.Store.Sessions.AuthenticationKey),
		[]byte(configs.Store.Sessions.EncryptionKey),
	)

	store.Options(sessions.Options{
		Path:     "/",
		Domain:   configs.Store.Sessions.Domain,
		MaxAge:   configs.Store.Sessions.MaxAge,
		Secure:   configs.Store.Sessions.Secure,
		HttpOnly: true,
	})
	sessionName = configs.Store.Sessions.Name
}

func Middleware() gin.HandlerFunc {
	return sessions.Sessions(sessionName, store)
}

func (sessionStruct *Session) Save(c *gin.Context) error {
	session := sessions.Default(c)
	session.Set(sessionName, sessionStruct)
	return session.Save()
}

func Get(c *gin.Context) *Session {
	sessionStruct, _ := sessions.Default(c).Get(sessionName).(Session)
	return &sessionStruct
}

func DeleteCookie(c *gin.Context) {
	c.SetCookie(sessionName,
		"",
		-1,
		"/",
		configs.Store.Sessions.Domain,
		configs.Store.Sessions.Secure,
		true,
	)
}
