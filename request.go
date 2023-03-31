package lightning

import (
	"io"
	"net/http"
)

type Request struct {
	req     *http.Request
	params  map[string]string
	method  string
	path    string
	RawBody []byte
}

// NewRequest creates a new Request object from an http.Request object
func NewRequest(req *http.Request) (*Request, error) {
	var rawBody []byte
	var err error
	if req.Body != nil {
		rawBody, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
	}

	request := &Request{
		req:     req,
		params:  map[string]string{},
		method:  req.Method,
		path:    req.URL.Path,
		RawBody: rawBody,
	}

	return request, nil
}

// SetParams sets the parameters for the Request object
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

// Queries returns the entire query parameter map for the context.
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

func (r *Request) UserAgent() string {
	return r.Header("user-agent")
}

func (r *Request) Referer() string {
	return r.Header("referer")
}

func (r *Request) RemoteAddr() string {
	ip := r.Header("x-real-ip")
	if ip == "" {
		ip = r.Header("x-forwarded-for")
		if ip == "" {
			ip = r.req.RemoteAddr
		}
	}
	return ip
}
