package lightning

import (
	"io"
	"net/http"
	"strings"
)

type request struct {
	req        *http.Request
	pathParams map[string]string
	method     string
	path       string
	rawBody    []byte
}

func newRequest(req *http.Request) (*request, error) {
	var rawBody []byte
	if req.Body != nil {
		var err error
		rawBody, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
	}

	return &request{
		req:        req,
		pathParams: map[string]string{},
		method:     req.Method,
		path:       req.URL.Path,
		rawBody:    rawBody,
	}, nil
}

func (r *request) setParams(params map[string]string) {
	r.pathParams = params
}

func (r *request) param(key string) string {
	return r.pathParams[key]
}

func (r *request) params() map[string]string {
	return r.pathParams
}

func (r *request) query(key string) string {
	return r.req.URL.Query().Get(key)
}

func (r *request) queries() map[string][]string {
	return r.req.URL.Query()
}

func (r *request) header(key string) string {
	return r.req.Header.Get(key)
}

func (r *request) headers() http.Header {
	return r.req.Header
}

func (r *request) cookie(name string) *http.Cookie {
	cookie, err := r.req.Cookie(name)
	if err != nil {
		return nil
	}
	return cookie
}

func (r *request) cookies() []*http.Cookie {
	return r.req.Cookies()
}

func (r *request) userAgent() string {
	return r.header("user-agent")
}

func (r *request) referer() string {
	return r.header("referer")
}

func (r *request) remoteAddr() string {
	ip := r.header("x-real-ip")
	if ip == "" {
		ip = r.header("x-forwarded-for")
		if ip != "" {
			if idx := strings.Index(ip, ","); idx != -1 {
				ip = strings.TrimSpace(ip[:idx])
			}
		}
	}
	if ip == "" {
		ip = r.req.RemoteAddr
	}
	return ip
}
