package lightning

import (
	"strings"

	"github.com/valyala/fasthttp"
)

type request struct {
	ctx        *fasthttp.RequestCtx
	pathParams map[string]string
}

func newRequest(ctx *fasthttp.RequestCtx) *request {
	return &request{
		ctx:        ctx,
		pathParams: make(map[string]string),
	}
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
	return string(r.ctx.QueryArgs().Peek(key))
}

func (r *request) queries() map[string][]string {
	args := r.ctx.QueryArgs()
	queries := make(map[string][]string)
	args.VisitAll(func(key, value []byte) {
		queries[string(key)] = append(queries[string(key)], string(value))
	})
	return queries
}

func (r *request) header(key string) string {
	return string(r.ctx.Request.Header.Peek(key))
}

func (r *request) headers() map[string]string {
	headers := make(map[string]string)
	r.ctx.Request.Header.VisitAll(func(key, value []byte) {
		headers[string(key)] = string(value)
	})
	return headers
}

func (r *request) cookie(name string) string {
	return string(r.ctx.Request.Header.Cookie(name))
}

func (r *request) cookies() map[string]string {
	cookies := make(map[string]string)
	r.ctx.Request.Header.VisitAllCookie(func(key, value []byte) {
		cookies[string(key)] = string(value)
	})
	return cookies
}

func (r *request) userAgent() string {
	return string(r.ctx.Request.Header.UserAgent())
}

func (r *request) referer() string {
	return string(r.ctx.Request.Header.Referer())
}

func (r *request) remoteAddr() string {
	ip := r.header("X-Real-IP")
	if ip == "" {
		ip = r.header("X-Forwarded-For")
		if ip != "" {
			if idx := strings.Index(ip, ","); idx != -1 {
				ip = strings.TrimSpace(ip[:idx])
			}
		}
	}
	if ip == "" {
		ip = r.ctx.RemoteAddr().String()
	}
	return ip
}

func (r *request) body() []byte {
	return r.ctx.Request.Body()
}

func (r *request) method() string {
	return string(r.ctx.Request.Header.Method())
}

func (r *request) path() string {
	return string(r.ctx.Path())
}

func (r *request) uri() string {
	return string(r.ctx.URI().FullURI())
}
