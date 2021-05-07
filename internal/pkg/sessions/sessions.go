/* Package sessions is ... */
package sessions

import (
	"encoding/gob"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

var (
	store cookie.Store
)

func Init() {
	gob.Register(time.Time{})

	store = cookie.NewStore(
		[]byte("4aceccc4ae3d43e28e7788c6165105e0"),
		[]byte("b9de910cded0409fb20655a7ccdc2b96"),
	)

	store.Options(sessions.Options{
		Path:     "/",
		Domain:   "localhost",
		MaxAge:   7776000,
		Secure:   false,
		HttpOnly: true,
	})
}

func Middleware() gin.HandlerFunc {
	return sessions.Sessions("session", store)
}
