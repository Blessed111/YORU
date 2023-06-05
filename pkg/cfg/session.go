package cfg

import (
	"github.com/gorilla/sessions"
)

const SESSION_ID = "user_logged_token"

var Store = sessions.NewCookieStore([]byte("ajksgk0934712qwfqqr"))

func init() {
	Store.Options.HttpOnly = true
	Store.Options.Secure = true
	Store.Options.MaxAge = 3600 * 24
}
