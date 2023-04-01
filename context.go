package lightning

import (
	"encoding/json"
	"net/http"
)

// Context represents the context of an HTTP request/response.
type Context struct {
	request  *Request
	response *response
	data     ContextData
	handlers []HandlerFunc
	index    int
	Method   string // HTTP method of the request
	Path     string // URL path of the request
}

// NewContext creates a new context object with the given HTTP response writer and request.
func NewContext(writer http.ResponseWriter, req *http.Request) (*Context, error) {
	request, err := NewRequest(req)
	if err != nil {
		return nil, err
	}
	response := newResponse(req, writer)
	ctx := &Context{
		request:  request,
		response: response,
		data:     ContextData{},
		handlers: []HandlerFunc{},
		index:    -1,
		Method:   request.method,
		Path:     request.path,
	}

	return ctx, nil
}

// Next calls the next middleware function in the chain.
func (c *Context) Next() {
	c.index++
	if c.index < len(c.handlers) {
		handlerFunc := c.handlers[c.index]
		handlerFunc(c)
	}
}

// Flush flushes the response buffer.
func (c *Context) Flush() {
	c.response.flush()
}

// RawBody returns the raw request body.
func (c *Context) RawBody() []byte {
	return c.request.RawBody
}

// StringBody returns the request body as a string.
func (c *Context) StringBody() string {
	return string(c.request.RawBody)
}

// JSONBody parses the request body as JSON and stores the result in v.
func (c *Context) JSONBody(v interface{}) error {
	err := json.Unmarshal(c.request.RawBody, v)
	if err != nil {
		return err
	}
	return nil
}

// SetHandlers sets the middleware handlers for the context.
func (c *Context) SetHandlers(handlers []HandlerFunc) {
	c.handlers = handlers
}

// SetParams sets the URL parameters for the request.
func (c *Context) SetParams(params map[string]string) {
	c.request.SetParams(params)
}

// Param returns the value of a URL parameter for a given key.
func (c *Context) Param(key string) string {
	return c.request.Param(key)
}

// Params returns all URL parameters for the request.
func (c *Context) Params() map[string]string {
	return c.request.Params()
}

// Query returns the value of a given query parameter.
func (c *Context) Query(key string) string {
	return c.request.Query(key)
}

// Queries returns all query parameters for the request.
func (c *Context) Queries() map[string][]string {
	return c.request.Queries()
}

// Status returns the HTTP status code of the response.
func (c *Context) Status() int {
	return c.response.status
}

// SetStatus sets the HTTP status code for the response.
func (c *Context) SetStatus(code int) {
	c.response.setStatus(code)
}

// Header returns the value of a given header.
func (c *Context) Header(key string) string {
	return c.request.Header(key)
}

// Headers returns all headers for the request.
func (c *Context) Headers() http.Header {
	return c.request.Headers()
}

// AddHeader adds a new header key-value pair to the response.
func (c *Context) AddHeader(key, value string) {
	c.response.addHeader(key, value)
}

// SetHeader sets the value of a given header in the response.
func (c *Context) SetHeader(key string, value string) {
	c.response.setHeader(key, value)
}

// DelHeader deletes a given header from the response.
func (c *Context) DelHeader(key string) {
	c.response.delHeader(key)
}

// Cookie returns the cookie with the given name.
func (c *Context) Cookie(name string) *http.Cookie {
	return c.request.Cookie(name)
}

// Cookies returns all cookies from the request.
func (c *Context) Cookies() []*http.Cookie {
	return c.request.Cookies()
}

// SetCookie sets a new cookie with the given key-value pair.
func (c *Context) SetCookie(key string, value string) {
	c.response.cookies.Set(key, value)
}

// SetCustomCookie sets a custom cookie in the response.
func (c *Context) SetCustomCookie(cookie *http.Cookie) {
	c.response.cookies.SetCustom(cookie)
}

// JSON writes a JSON response with the given status code and object.
func (c *Context) JSON(code int, obj interface{}) {
	c.response.setStatus(code)
	err := c.response.json(obj)
	if err != nil {
		panic(err)
	}
}

// Text writes a plain text response with the given status code and format.
func (c *Context) Text(code int, text string) {
	c.response.setStatus(code)
	c.response.text(text)
}

// XML writes an XML response with the given status code and object.
func (c *Context) XML(code int, obj interface{}) {
	c.response.setStatus(code)
	err := c.response.xml(obj)
	if err != nil {
		return
	}
}

// File writes a file as the response.
func (c *Context) File(filepath string) {
	err := c.response.file(filepath)
	if err != nil {
		return
	}
}

// GetData returns the value of a custom data field for the context.
func (c *Context) GetData(key string) interface{} {
	return c.data.Get(key)
}

// SetData sets the value of a custom data field for the context.
func (c *Context) SetData(key string, value interface{}) {
	c.data.Set(key, value)
}

// DelData deletes a custom data field from the context.
func (c *Context) DelData(key string) {
	c.data.Del(key)
}

// Redirect redirects the request to a new URL with the given status code.
func (c *Context) Redirect(code int, url string) {
	c.response.redirect(code, url)
}

func (c *Context) UserAgent() string {
	return c.request.UserAgent()
}

func (c *Context) Referer() string {
	return c.request.Referer()
}

func (c *Context) RemoteAddr() string {
	return c.request.RemoteAddr()
}

func (c *Context) Success(data interface{}) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"msg":  "ok",
		"data": data,
	})
}

func (c *Context) Fail(code int, msg string) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"code": code,
		"msg":  msg,
	})
}
