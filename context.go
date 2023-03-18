package lightning

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Context struct {
	Response http.ResponseWriter
	Request  *http.Request
	Method   string
	Path     string
	params   map[string]string
	cookies  map[string]*http.Cookie
}

// NewContext creates a new context object with the given HTTP response writer and request.
func NewContext(writer http.ResponseWriter, req *http.Request, params map[string]string) *Context {
	return &Context{
		Response: writer,
		Request:  req,
		Method:   req.Method,
		Path:     req.URL.Path,
		params:   params,
	}
}

// Param returns the parameter value for a given key.
func (c *Context) Param(key string) string {
	return c.params[key]
}

// Params returns the entire parameter map for the context.
func (c *Context) Params() map[string]string {
	return c.params
}

// Query returns the value of a given query parameter.
func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

// SetStatus sets the HTTP status code for the response.
func (c *Context) SetStatus(code int) {
	c.Response.WriteHeader(code)
}

// Header returns the value of a given header.
func (c *Context) Header(key string) string {
	return c.Request.Header.Get(key)
}

// Headers returns the entire header map for the request.
func (c *Context) Headers() http.Header {
	return c.Request.Header
}

// AddHeader adds a new header key-value pair to the response.
func (c *Context) AddHeader(key, value string) {
	c.Response.Header().Add(key, value)
}

// SetHeader sets the value of a given header in the response.
func (c *Context) SetHeader(key string, value string) {
	c.Response.Header().Set(key, value)
}

// DelHeader deletes a given header from the response.
func (c *Context) DelHeader(key string) {
	c.Response.Header().Del(key)
}

// Cookie returns the cookie with the given name.
func (c *Context) Cookie(name string) (*http.Cookie, error) {
	return c.Request.Cookie(name)
}

// Cookies returns all cookies from the request.
func (c *Context) Cookies() []*http.Cookie {
	return c.Request.Cookies()
}

// SetCookie sets a new cookie with the given key-value pair.
func (c *Context) SetCookie(key string, value string) {
	cookie := &http.Cookie{Name: key, Value: value, Path: "/"}
	http.SetCookie(c.Response, cookie)
}

// SetCustomCookie sets a custom cookie in the response.
func (c *Context) SetCustomCookie(cookie *http.Cookie) {
	http.SetCookie(c.Response, cookie)
}

// Text writes a plain text response with the given status code and format.
func (c *Context) Text(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.SetStatus(code)
	_, err := fmt.Fprintf(c.Response, format, values...)
	if err != nil {
		return
	}
}

// JSON writes a JSON response with the given status code and object.
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.SetStatus(code)
	encoder := json.NewEncoder(c.Response)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Response, err.Error(), 500)
	}
}
