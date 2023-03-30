package lightning

import (
	"net/http"
)

type Context struct {
	request  *Request
	response *Response
	data     ContextData
	handlers []HandlerFunc
	index    int
	Method   string
	Path     string
}

// NewContext creates a new context object with the given HTTP response writer and request.
func NewContext(writer http.ResponseWriter, req *http.Request) *Context {
	request := NewRequest(req)
	response := NewResponse(req, writer)
	ctx := &Context{
		request:  request,
		response: response,
		data:     ContextData{},
		handlers: []HandlerFunc{},
		index:    -1,
		Method:   request.method,
		Path:     request.path,
	}

	return ctx
}

func (c *Context) Next() {
	c.index++
	if c.index < len(c.handlers) {
		handlerFunc := c.handlers[c.index]
		handlerFunc(c)
	}
}

func (c *Context) Flush() {
	c.response.flush()
}

func (c *Context) SetHandlers(handlers []HandlerFunc) {
	c.handlers = handlers
}

func (c *Context) SetParams(params map[string]string) {
	c.request.SetParams(params)
}

// Param returns the parameter value for a given key.
func (c *Context) Param(key string) string {
	return c.request.Param(key)
}

// Params returns the entire parameter map for the context.
func (c *Context) Params() map[string]string {
	return c.request.Params()
}

// Query returns the value of a given query parameter.
func (c *Context) Query(key string) string {
	return c.request.Query(key)
}

func (c *Context) Queries() map[string][]string {
	return c.request.Queries()
}

func (c *Context) Status() int {
	return c.response.status
}

// SetStatus sets the HTTP status code for the response.
func (c *Context) SetStatus(code int) {
	c.response.SetStatus(code)
}

// Header returns the value of a given header.
func (c *Context) Header(key string) string {
	return c.request.Header(key)
}

// Headers returns the entire header map for the request.
func (c *Context) Headers() http.Header {
	return c.request.Headers()
}

// AddHeader adds a new header key-value pair to the response.
func (c *Context) AddHeader(key, value string) {
	c.response.AddHeader(key, value)
}

// SetHeader sets the value of a given header in the response.
func (c *Context) SetHeader(key string, value string) {
	c.response.SetHeader(key, value)
}

// DelHeader deletes a given header from the response.
func (c *Context) DelHeader(key string) {
	c.response.DelHeader(key)
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
	c.response.Cookies.Set(key, value)
}

// SetCustomCookie sets a custom cookie in the response.
func (c *Context) SetCustomCookie(cookie *http.Cookie) {
	c.response.Cookies.SetCustom(cookie)
}

// JSON writes a JSON response with the given status code and object.
func (c *Context) JSON(obj interface{}) {
	c.response.SetStatus(http.StatusOK)
	err := c.response.JSON(obj)
	if err != nil {
		return
	}
}

// Text writes a plain text response with the given status code and format.
func (c *Context) Text(text string) {
	c.response.SetStatus(http.StatusOK)
	c.response.Text(text)
}

func (c *Context) XML(obj interface{}) {
	c.response.SetStatus(http.StatusOK)
	err := c.response.XML(obj)
	if err != nil {
		return
	}
}

func (c *Context) File(filepath string) {
	err := c.response.File(filepath)
	if err != nil {
		return
	}
}

func (c *Context) NotFound() {
	c.response.SetStatus(http.StatusNotFound)
	c.response.Text(http.StatusText(http.StatusNotFound))
}

func (c *Context) GetData(key string) interface{} {
	return c.data.Get(key)
}

func (c *Context) SetData(key string, value interface{}) {
	c.data.Set(key, value)
}

func (c *Context) DelData(key string) {
	c.data.Del(key)
}
