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
}

func Middleware() gin.HandlerFunc {
	return sessions.Sessions(configs.Store.Sessions.Name, store)
}

func (sessionStruct *Session) Save(c *gin.Context) error {
	session := sessions.Default(c)
	session.Set(configs.Store.Sessions.Name, sessionStruct)
	return session.Save()
}

func Get(c *gin.Context) *Session {
	sessionStruct, _ := sessions.Default(c).Get(configs.Store.Sessions.Name).(Session)
	return &sessionStruct
}

func DeleteCookie(c *gin.Context) {
	c.SetCookie(configs.Store.Sessions.Name,
		"",
		-1,
		"/",
		configs.Store.Sessions.Domain,
		configs.Store.Sessions.Secure,
		true,
	)
}
