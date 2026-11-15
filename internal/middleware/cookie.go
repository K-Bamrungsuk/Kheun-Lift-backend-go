package middleware

import (
	"net/http"
	"time"
)

const sessionCookieName = "kheun_session"

// SessionCookie holds the settings shared by login, logout and delete account.
type SessionCookie struct {
	Secure bool
	Domain string // empty locally; ".<domain>" on prod so app.<domain> also gets the cookie
}

func (s SessionCookie) Set(w http.ResponseWriter, value string, expires time.Time) {
	s.write(w, value, 86400, expires)
}

func (s SessionCookie) Clear(w http.ResponseWriter) {
	s.write(w, "", -1, time.Unix(0, 0).UTC())
}

// Browsers delete a cookie only when Domain and Path match the ones it was set with.
func (s SessionCookie) write(w http.ResponseWriter, value string, maxAge int, expires time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    value,
		Path:     "/",
		Domain:   s.Domain,
		HttpOnly: true,
		Secure:   s.Secure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   maxAge,
		Expires:  expires,
	})
}
