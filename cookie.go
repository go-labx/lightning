package lightning

import (
	"net/http"
)

// cookiesMap is a map of http.Cookie pointers.
type cookiesMap map[string]*http.Cookie

// get returns the http.Cookie pointer with the given key.
func (cookies cookiesMap) get(key string) *http.Cookie {
	return cookies[key]
}

// set sets the http.Cookie with the given key and value.
func (cookies cookiesMap) set(key string, value string) {
	cookies[key] = &http.Cookie{
		Name:  key,
		Value: value,
		Path:  "/",
	}
}

// del deletes the http.Cookie with the given key.
func (cookies cookiesMap) del(key string) {
	delete(cookies, key)
}

// setCustom sets the given http.Cookie pointer.
func (cookies cookiesMap) setCustom(cookie *http.Cookie) {
	cookies[cookie.Name] = cookie
}
