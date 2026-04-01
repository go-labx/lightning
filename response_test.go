package lightning

import (
	"testing"

	"github.com/valyala/fasthttp"
)

func createResponse() (*response, *fasthttp.RequestCtx) {
	ctx := &fasthttp.RequestCtx{}
	resp := newResponse(ctx)
	resp.cookies = make(cookiesMap)
	return resp, ctx
}

func TestNewResponse(t *testing.T) {
	resp, ctx := createResponse()
	if resp.ctx != ctx {
		t.Error("ctx not set correctly")
	}
	if resp.statusCode != StatusNotFound {
		t.Errorf("Expected default status %d, got %d", StatusNotFound, resp.statusCode)
	}
}

func TestResponse_setStatus(t *testing.T) {
	resp, _ := createResponse()

	resp.setStatus(StatusOK)
	if resp.statusCode != StatusOK {
		t.Errorf("Expected status %d, got %d", StatusOK, resp.statusCode)
	}
}

func TestResponse_setBody(t *testing.T) {
	resp, _ := createResponse()

	resp.setBody([]byte("test body"))
	if string(resp.body) != "test body" {
		t.Errorf("Expected body 'test body', got '%s'", string(resp.body))
	}
}

func TestResponse_redirect(t *testing.T) {
	resp, _ := createResponse()

	resp.redirect(StatusMovedPermanently, "https://example.com")
	if resp.redirectTo != "https://example.com" {
		t.Errorf("Expected redirect URL 'https://example.com', got '%s'", resp.redirectTo)
	}
	if resp.statusCode != StatusMovedPermanently {
		t.Errorf("Expected status %d, got %d", StatusMovedPermanently, resp.statusCode)
	}
}

func TestResponse_addHeader(t *testing.T) {
	resp, ctx := createResponse()

	resp.addHeader("X-Custom", "value1")
	resp.addHeader("X-Custom", "value2")

	hdr := string(ctx.Response.Header.Peek("X-Custom"))
	if hdr == "" {
		t.Error("Expected header to be set")
	}
}

func TestResponse_setHeader(t *testing.T) {
	resp, ctx := createResponse()

	resp.setHeader("Content-Type", "application/json")

	hdr := string(ctx.Response.Header.Peek("Content-Type"))
	if hdr != "application/json" {
		t.Errorf("Expected 'application/json', got '%s'", hdr)
	}
}

func TestResponse_delHeader(t *testing.T) {
	resp, ctx := createResponse()

	resp.setHeader("X-Custom", "value")
	resp.delHeader("X-Custom")

	hdr := string(ctx.Response.Header.Peek("X-Custom"))
	if hdr != "" {
		t.Errorf("Expected empty header, got '%s'", hdr)
	}
}

func TestResponse_file(t *testing.T) {
	resp, _ := createResponse()

	resp.file("/nonexistent/path")
	resp.flush()
}

func TestResponse_flush(t *testing.T) {
	resp, ctx := createResponse()

	resp.setStatus(StatusOK)
	resp.setBody([]byte("test body"))
	resp.flush()

	if ctx.Response.StatusCode() != StatusOK {
		t.Errorf("Expected status %d, got %d", StatusOK, ctx.Response.StatusCode())
	}
	if string(ctx.Response.Body()) != "test body" {
		t.Errorf("Expected body 'test body', got '%s'", string(ctx.Response.Body()))
	}
}

func TestResponse_flushWithCookie(t *testing.T) {
	resp, ctx := createResponse()

	resp.cookies.set("session", "abc123")
	resp.flush()

	cookie := string(ctx.Response.Header.Peek("Set-Cookie"))
	if cookie == "" {
		t.Error("Expected Set-Cookie header to be set")
	}
}

func TestResponse_flushWithRedirect(t *testing.T) {
	resp, ctx := createResponse()

	resp.redirect(StatusFound, "/new-location")
	resp.flush()

	if ctx.Response.StatusCode() != StatusFound {
		t.Errorf("Expected status %d, got %d", StatusFound, ctx.Response.StatusCode())
	}
	location := string(ctx.Response.Header.Peek("Location"))
	if location == "" {
		t.Error("Expected Location header to be set")
	}
}
