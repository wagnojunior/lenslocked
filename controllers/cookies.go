package controllers

import (
	"fmt"
	"net/http"
)

const (
	CookieSession = "session"
)

// newCookie returns a new cookie with with fixed path and http only.
func newCookie(name, value string) *http.Cookie {
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
	}
	return &cookie
}

// setCookie sets a new cookie to the response writer
func setCookie(w http.ResponseWriter, name, value string) {
	http.SetCookie(w, newCookie(name, value))
}

// readCookie reads the cookie defined by `name` and returns its value and an error
func readCookie(r *http.Request, name string) (string, error) {
	c, err := r.Cookie(name)
	if err != nil {
		return "", fmt.Errorf("%s: %w", name, err)
	}
	return c.Value, nil
}

// deleteCookie deletes a cookie by overwriting the existing one with a negative `MaxAge` field
func deleteCookie(w http.ResponseWriter, name string) {
	cookie := newCookie(name, "")
	cookie.MaxAge = -1
	http.SetCookie(w, cookie)
}
