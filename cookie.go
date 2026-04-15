package lightning

import (
	"github.com/valyala/fasthttp"
)

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

// CookieConfig holds the configuration for setting a cookie with security attributes.
type CookieConfig struct {
	Name     string
	Value    string
	Path     string
	Domain   string
	MaxAge   int
	Secure   bool
	HttpOnly bool
	SameSite string
}

// isCookieSameSiteSupported is a sentinel to check fasthttp SameSite availability.
var isCookieSameSiteSupported = true

// setCookieWithConfig applies a CookieConfig to a fasthttp.Cookie.
func setCookieWithConfig(fc *fasthttp.Cookie, config CookieConfig) {
	fc.SetKey(config.Name)
	fc.SetValue(config.Value)
	if config.Path != "" {
		fc.SetPath(config.Path)
	}
	if config.Domain != "" {
		fc.SetDomain(config.Domain)
	}
	if config.MaxAge > 0 {
		fc.SetMaxAge(config.MaxAge)
	}
	fc.SetSecure(config.Secure)
	fc.SetHTTPOnly(config.HttpOnly)
	switch config.SameSite {
	case "Strict":
		fc.SetSameSite(fasthttp.CookieSameSiteStrictMode)
	case "Lax":
		fc.SetSameSite(fasthttp.CookieSameSiteLaxMode)
	case "None":
		fc.SetSameSite(fasthttp.CookieSameSiteNoneMode)
	}
}
