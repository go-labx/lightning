package lightning

import (
	"io"
	"net/http"
)

type request struct {
	originReq *http.Request
	paramsMap map[string]string
	method    string
	path      string
	rawBody   []byte
}

// newRequest creates a new request object from an http.Request object
func newRequest(req *http.Request) (*request, error) {
	var rawBody []byte
	var err error
	if req.Body != nil {
		rawBody, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
	}

	request := &request{
		originReq: req,
		paramsMap: map[string]string{},
		method:    req.Method,
		path:      req.URL.Path,
		rawBody:   rawBody,
	}

	return request, nil
}

// setParams sets the parameters for the request object
func (r *request) setParams(params map[string]string) {
	r.paramsMap = params
}

// param returns the parameter value for a given key.
func (r *request) param(key string) string {
	return r.paramsMap[key]
}

// params returns the entire parameter map for the context.
func (r *request) params() map[string]string {
	return r.paramsMap
}

// query returns the value of a given query parameter.
func (r *request) query(key string) string {
	return r.originReq.URL.Query().Get(key)
}

// queries returns the entire query parameter map for the context.
func (r *request) queries() map[string][]string {
	return r.originReq.URL.Query()
}

// header returns the value of a given header.
func (r *request) header(key string) string {
	return r.originReq.Header.Get(key)
}

// headers returns the entire header map for the request.
func (r *request) headers() http.Header {
	return r.originReq.Header
}

// cookie returns the cookie with the given name.
func (r *request) cookie(name string) *http.Cookie {
	cookie, err := r.originReq.Cookie(name)
	if err != nil {
		return nil
	}
	return cookie
}

// cookiesMap returns all cookies from the request.
func (r *request) cookies() []*http.Cookie {
	return r.originReq.Cookies()
}

// userAgent returns the user agent header value of the request.
func (r *request) userAgent() string {
	return r.header("user-agent")
}

// referer returns the referer header value of the request.
func (r *request) referer() string {
	return r.header("referer")
}

// remoteAddr returns the remote address of the request.
func (r *request) remoteAddr() string {
	ip := r.header("x-real-ip")
	if ip == "" {
		ip = r.header("x-forwarded-for")
		if ip == "" {
			ip = r.originReq.RemoteAddr
		}
	}
	return ip
}
