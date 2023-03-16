package lightning

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Context struct {
	writer http.ResponseWriter
	req    *http.Request
	Method string
	Path   string
}

func NewContext(writer http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		writer: writer,
		req:    req,
		Method: req.Method,
		Path:   req.URL.Path,
	}
}

func (c *Context) Param(key string) string {
	return c.req.URL.Query().Get(key)
}

func (c *Context) SetStatus(code int) {
	c.writer.WriteHeader(code)
}

func (c *Context) SetHeader(key string, value string) {
	c.writer.Header().Set(key, value)
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.SetStatus(code)
	_, err := fmt.Fprintf(c.writer, format, values...)
	if err != nil {
		return
	}
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.SetStatus(code)
	encoder := json.NewEncoder(c.writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.writer, err.Error(), 500)
	}
}
