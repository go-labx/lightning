package lightning

import "net/http"

type Request struct {
	req    *http.Request
	params map[string]string
	method string
	path   string
}

func NewRequest(req *http.Request) *Request {
	request := &Request{
		req:    req,
		params: map[string]string{},
		method: req.Method,
		path:   req.URL.Path,
	}

	return request
}

func (r *Request) SetParams(params map[string]string) {
	r.params = params
}

// Param returns the parameter value for a given key.
func (r *Request) Param(key string) string {
	return r.params[key]
}

// Params returns the entire parameter map for the context.
func (r *Request) Params() map[string]string {
	return r.params
}

// Query returns the value of a given query parameter.
func (r *Request) Query(key string) string {
	return r.req.URL.Query().Get(key)
}

func (r *Request) Queries() map[string][]string {
	return r.req.URL.Query()
}

// Header returns the value of a given header.
func (r *Request) Header(key string) string {
	return r.req.Header.Get(key)
}

// Headers returns the entire header map for the request.
func (r *Request) Headers() http.Header {
	return r.req.Header
}

// Cookie returns the cookie with the given name.
func (r *Request) Cookie(name string) *http.Cookie {
	cookie, err := r.req.Cookie(name)
	if err != nil {
		return nil
	}
	return cookie
}

// Cookies returns all cookies from the request.
func (r *Request) Cookies() []*http.Cookie {
	return r.req.Cookies()
}
