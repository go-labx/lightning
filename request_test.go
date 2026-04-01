package lightning

import (
	"testing"

	"github.com/valyala/fasthttp"
)

func createFasthttpCtx(method, path string) (*fasthttp.RequestCtx, *request) {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod(method)
	ctx.Request.Header.SetRequestURI(path)
	r := newRequest(ctx)
	return ctx, r
}

func TestRequest_Cookie(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod("GET")
	ctx.Request.Header.SetRequestURI("/test")
	ctx.Request.Header.SetCookie("cookie1", "value1")

	r := newRequest(ctx)

	cookie := r.cookie("cookie1")
	if cookie != nil && string(cookie.Key()) != "cookie1" {
		t.Errorf("Expected cookie key 'cookie1', got '%s'", string(cookie.Key()))
	}

	cookie = r.cookie("nonexistent")
	if cookie != nil {
		t.Errorf("Expected nil for nonexistent cookie, got %v", cookie)
	}
}

func TestRequest_Cookies(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod("GET")
	ctx.Request.Header.SetRequestURI("/test")
	ctx.Request.Header.SetCookie("cookie1", "value1")
	ctx.Request.Header.SetCookie("cookie2", "value2")

	r := newRequest(ctx)
	cookies := r.cookies()

	if len(cookies) == 0 {
		t.Error("Expected cookies, got empty")
	}
}

func TestRequest_Header(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod("GET")
	ctx.Request.Header.SetRequestURI("/test")
	ctx.Request.Header.Set("header1", "value1")

	r := newRequest(ctx)

	if got := r.header("header1"); got != "value1" {
		t.Errorf("header() = %v, want %v", got, "value1")
	}

	if got := r.header("nonexistent"); got != "" {
		t.Errorf("header() = %v, want empty", got)
	}
}

func TestRequest_Headers(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod("GET")
	ctx.Request.Header.SetRequestURI("/test")
	ctx.Request.Header.Set("Content-Type", "application/json")

	r := newRequest(ctx)
	headers := r.headers()

	if len(headers) == 0 {
		t.Error("Expected headers to be non-empty")
	}
}

func TestRequest_Param(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod("GET")
	ctx.Request.Header.SetRequestURI("/test")

	r := newRequest(ctx)
	r.setParams(map[string]string{"param1": "value1"})

	if got := r.param("param1"); got != "value1" {
		t.Errorf("param() = %v, want %v", got, "value1")
	}

	if got := r.param("nonexistent"); got != "" {
		t.Errorf("param() = %v, want empty", got)
	}
}

func TestRequest_Params(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod("GET")
	ctx.Request.Header.SetRequestURI("/test")

	r := newRequest(ctx)
	params := map[string]string{"param1": "value1", "param2": "value2"}
	r.setParams(params)

	got := r.params()
	for k, v := range params {
		if got[k] != v {
			t.Errorf("params()[%s] = %v, want %v", k, got[k], v)
		}
	}
}

func TestRequest_Queries(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod("GET")
	ctx.Request.Header.SetRequestURI("/test?param1=value1&param2=value2")

	r := newRequest(ctx)
	queries := r.queries()

	if queries["param1"][0] != "value1" {
		t.Errorf("Expected param1=value1, got %v", queries["param1"])
	}
	if queries["param2"][0] != "value2" {
		t.Errorf("Expected param2=value2, got %v", queries["param2"])
	}
}

func TestRequest_Query(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod("GET")
	ctx.Request.Header.SetRequestURI("/test?param1=value1")

	r := newRequest(ctx)

	if got := r.query("param1"); got != "value1" {
		t.Errorf("query() = %v, want %v", got, "value1")
	}

	if got := r.query("nonexistent"); got != "" {
		t.Errorf("query() = %v, want empty", got)
	}
}

func TestRequest_userAgent(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod("GET")
	ctx.Request.Header.SetRequestURI("/test")
	ctx.Request.Header.SetUserAgent("test-user-agent")

	r := newRequest(ctx)

	if got := r.userAgent(); got != "test-user-agent" {
		t.Errorf("userAgent() = %v, want %v", got, "test-user-agent")
	}
}

func TestRequest_referer(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod("GET")
	ctx.Request.Header.SetRequestURI("/test")
	ctx.Request.Header.SetReferer("https://example.com")

	r := newRequest(ctx)

	if got := r.referer(); got != "https://example.com" {
		t.Errorf("referer() = %v, want %v", got, "https://example.com")
	}
}

func TestRequest_remoteAddr(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod("GET")
	ctx.Request.Header.SetRequestURI("/test")

	r := newRequest(ctx)
	addr := r.remoteAddr()

	if addr == "" {
		t.Error("Expected non-empty remote address")
	}
}

func TestRequest_remoteAddrWithXRealIP(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod("GET")
	ctx.Request.Header.SetRequestURI("/test")
	ctx.Request.Header.Set("X-Real-IP", "1.2.3.4")

	r := newRequest(ctx)
	addr := r.remoteAddr()

	if addr != "1.2.3.4" {
		t.Errorf("Expected X-Real-IP, got %s", addr)
	}
}

func TestRequest_remoteAddrWithXForwardedFor(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod("GET")
	ctx.Request.Header.SetRequestURI("/test")
	ctx.Request.Header.Set("X-Forwarded-For", "1.2.3.4, 5.6.7.8")

	r := newRequest(ctx)
	addr := r.remoteAddr()

	if addr != "1.2.3.4" {
		t.Errorf("Expected first IP from X-Forwarded-For, got %s", addr)
	}
}

func TestRequest_body(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod("POST")
	ctx.Request.Header.SetRequestURI("/test")
	ctx.Request.SetBody([]byte("test body"))

	r := newRequest(ctx)
	body := r.body()

	if string(body) != "test body" {
		t.Errorf("body() = %v, want %v", string(body), "test body")
	}
}

func TestRequest_method(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod("POST")
	ctx.Request.Header.SetRequestURI("/test")

	r := newRequest(ctx)

	if got := r.method(); got != "POST" {
		t.Errorf("method() = %v, want %v", got, "POST")
	}
}

func TestRequest_path(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod("GET")
	ctx.Request.Header.SetRequestURI("/test/path")

	r := newRequest(ctx)

	if got := r.path(); got != "/test/path" {
		t.Errorf("path() = %v, want %v", got, "/test/path")
	}
}

func TestRequest_uri(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod("GET")
	ctx.Request.Header.SetRequestURI("/test?param=value")

	r := newRequest(ctx)
	uri := r.uri()

	if uri == "" {
		t.Error("Expected non-empty URI")
	}
}
