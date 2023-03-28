package lightning

import (
	"net/http"
)

// Cookie is a map of http.Cookie pointers.
type Cookie map[string]*http.Cookie

// Get returns the http.Cookie pointer with the given key.
func (cookies Cookie) Get(key string) *http.Cookie {
	return cookies[key]
}

// Set sets the http.Cookie with the given key and value.
func (cookies Cookie) Set(key string, value string) {
	cookies[key] = &http.Cookie{
		Name:  key,
		Value: value,
		Path:  "/",
	}
}

// Del deletes the http.Cookie with the given key.
func (cookies Cookie) Del(key string) {
	delete(cookies, key)
}

// SetCustom sets the given http.Cookie pointer.
func (cookies Cookie) SetCustom(cookie *http.Cookie) {
	cookies[cookie.Name] = cookie
}
