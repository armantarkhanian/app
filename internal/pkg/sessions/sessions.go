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
		[]byte("4aceccc4ae3d43e28e7788c6165105e0"),
		[]byte("b9de910cded0409fb20655a7ccdc2b96"),
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
	return sessions.Sessions("session", store)
}

func (sessionStruct *Session) Save(c *gin.Context) error {
	session := sessions.Default(c)
	session.Set("session", sessionStruct)
	return session.Save()
}

func Get(c *gin.Context) *Session {
	sessionStruct, _ := sessions.Default(c).Get("session").(Session)
	return &sessionStruct
}

func DeleteCookie(c *gin.Context) {
	c.SetCookie("session",
		"",
		-1,
		"/",
		configs.Store.Sessions.Domain,
		configs.Store.Sessions.Secure,
		true,
	)
}
