package lightning

// cookiesMap is a map of cookie values keyed by cookie name.
type cookiesMap map[string]string

// set sets the cookie value with the given key.
func (cookies cookiesMap) set(key string, value string) {
	cookies[key] = value
}

// del deletes the cookie with the given key.
func (cookies cookiesMap) del(key string) {
	delete(cookies, key)
}
