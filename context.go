package lightning

import (
	"encoding/json"
	"net/http"
)

// Context represents the context of an HTTP request/response.
type Context struct {
	req      *request
	res      *response
	data     contextData
	handlers []HandlerFunc
	index    int
	Method   string // HTTP method of the originReq
	Path     string // URL path of the originReq
}

// newContext creates a new context object with the given HTTP response writer and req.
func newContext(writer http.ResponseWriter, req *http.Request) (*Context, error) {
	request, err := newRequest(req)
	if err != nil {
		return nil, err
	}
	response := newResponse(req, writer)
	ctx := &Context{
		req:      request,
		res:      response,
		data:     contextData{},
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

// flushResponse flushes the response buffer.
func (c *Context) flushResponse() {
	c.res.flush()
}

// RawBody returns the raw origin request body.
func (c *Context) RawBody() []byte {
	return c.req.rawBody
}

// StringBody returns the origin request body as a string.
func (c *Context) StringBody() string {
	return string(c.req.rawBody)
}

// JSONBody parses the origin request body as JSON and stores the result in v.
func (c *Context) JSONBody(v interface{}) error {
	err := json.Unmarshal(c.req.rawBody, v)
	if err != nil {
		return err
	}
	return nil
}

// setHandlers sets the handlers for the context.
func (c *Context) setHandlers(handlers []HandlerFunc) {
	c.handlers = handlers
}

// setParams sets the URL parameters for the req.
func (c *Context) setParams(params map[string]string) {
	c.req.setParams(params)
}

// Param returns the value of a URL parameter for a given key.
func (c *Context) Param(key string) string {
	return c.req.param(key)
}

// Params returns all URL parameters for the req.
func (c *Context) Params() map[string]string {
	return c.req.params()
}

// Query returns the value of a given query parameter.
func (c *Context) Query(key string) string {
	return c.req.query(key)
}

// Queries returns all query parameters for the req.
func (c *Context) Queries() map[string][]string {
	return c.req.queries()
}

// Status returns the HTTP status code of the response.
func (c *Context) Status() int {
	return c.res.statusCode
}

// SetStatus sets the HTTP status code for the response.
func (c *Context) SetStatus(code int) {
	c.res.setStatus(code)
}

// Header returns the value of a given header.
func (c *Context) Header(key string) string {
	return c.req.header(key)
}

// Headers returns all headers for the req.
func (c *Context) Headers() http.Header {
	return c.req.headers()
}

// AddHeader adds a new header key-value pair to the response.
func (c *Context) AddHeader(key, value string) {
	c.res.addHeader(key, value)
}

// SetHeader sets the value of a given header in the response.
func (c *Context) SetHeader(key string, value string) {
	c.res.setHeader(key, value)
}

// DelHeader deletes a given header from the response.
func (c *Context) DelHeader(key string) {
	c.res.delHeader(key)
}

// Cookie returns the cookie with the given name.
func (c *Context) Cookie(name string) *http.Cookie {
	return c.req.cookie(name)
}

// Cookies returns all cookies from the req.
func (c *Context) Cookies() []*http.Cookie {
	return c.req.cookies()
}

// SetCookie sets a new cookie with the given key-value pair.
func (c *Context) SetCookie(key string, value string) {
	c.res.cookies.set(key, value)
}

// SetCustomCookie sets a custom cookie in the response.
func (c *Context) SetCustomCookie(cookie *http.Cookie) {
	c.res.cookies.setCustom(cookie)
}

// JSON writes a JSON response with the given status code and object.
func (c *Context) JSON(code int, obj interface{}) {
	c.res.setStatus(code)
	err := c.res.json(obj)
	if err != nil {
		panic(err)
	}
}

// Text writes a plain text response with the given status code and format.
func (c *Context) Text(code int, text string) {
	c.res.setStatus(code)
	c.res.text(text)
}

// XML writes an XML response with the given status code and object.
func (c *Context) XML(code int, obj interface{}) {
	c.res.setStatus(code)
	err := c.res.xml(obj)
	if err != nil {
		return
	}
}

// File writes a file as the response.
func (c *Context) File(filepath string) {
	err := c.res.file(filepath)
	if err != nil {
		return
	}
}

// GetData returns the value of a custom data field for the context.
func (c *Context) GetData(key string) interface{} {
	return c.data.get(key)
}

// SetData sets the value of a custom data field for the context.
func (c *Context) SetData(key string, value interface{}) {
	c.data.set(key, value)
}

// DelData deletes a custom data field from the context.
func (c *Context) DelData(key string) {
	c.data.del(key)
}

// Redirect redirects the originReq to a new URL with the given status code.
func (c *Context) Redirect(code int, url string) {
	c.res.redirect(code, url)
}

func (c *Context) UserAgent() string {
	return c.req.userAgent()
}

func (c *Context) Referer() string {
	return c.req.referer()
}

func (c *Context) RemoteAddr() string {
	return c.req.remoteAddr()
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
