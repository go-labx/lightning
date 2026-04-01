package lightning

import (
	"testing"

	"github.com/valyala/fasthttp"
)

func newTestCtx(method, path string) *fasthttp.RequestCtx {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod(method)
	ctx.Request.Header.SetRequestURI(path)
	return ctx
}

func TestIntegrationBasicRouting(t *testing.T) {
	app := NewApp()

	app.Get("/get", func(c *Context) {
		c.Text(StatusOK, "GET response")
	})

	app.Post("/post", func(c *Context) {
		c.Text(StatusOK, "POST response")
	})

	app.Put("/put", func(c *Context) {
		c.Text(StatusOK, "PUT response")
	})

	app.Delete("/delete", func(c *Context) {
		c.Text(StatusOK, "DELETE response")
	})

	tests := []struct {
		method string
		path   string
		want   int
	}{
		{MethodGet, "/get", StatusOK},
		{MethodPost, "/post", StatusOK},
		{MethodPut, "/put", StatusOK},
		{MethodDelete, "/delete", StatusOK},
		{MethodGet, "/nonexistent", StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.method+" "+tt.path, func(t *testing.T) {
			ctx := newTestCtx(tt.method, tt.path)
			app.serveRequest(ctx)
			if ctx.Response.StatusCode() != tt.want {
				t.Errorf("Expected status %d, got %d", tt.want, ctx.Response.StatusCode())
			}
		})
	}
}

func TestIntegrationJSONAPI(t *testing.T) {
	app := NewApp()

	type Response struct {
		Message string `json:"message"`
	}

	app.Get("/api/json", func(c *Context) {
		c.JSON(StatusOK, Response{Message: "hello"})
	})

	ctx := newTestCtx(MethodGet, "/api/json")
	app.serveRequest(ctx)

	if ctx.Response.StatusCode() != StatusOK {
		t.Errorf("Expected status %d, got %d", StatusOK, ctx.Response.StatusCode())
	}
}

func TestIntegrationMiddlewareChain(t *testing.T) {
	app := NewApp()
	called := []string{}

	app.Use(func(c *Context) {
		called = append(called, "middleware1-before")
		c.Next()
		called = append(called, "middleware1-after")
	})

	app.Use(func(c *Context) {
		called = append(called, "middleware2-before")
		c.Next()
		called = append(called, "middleware2-after")
	})

	app.Get("/test", func(c *Context) {
		called = append(called, "handler")
	})

	ctx := newTestCtx(MethodGet, "/test")
	app.serveRequest(ctx)

	expected := []string{
		"middleware1-before",
		"middleware2-before",
		"handler",
		"middleware2-after",
		"middleware1-after",
	}

	if len(called) != len(expected) {
		t.Fatalf("Expected %d calls, got %d: %v", len(expected), len(called), called)
	}

	for i, e := range expected {
		if called[i] != e {
			t.Errorf("Call[%d] = %s, want %s", i, called[i], e)
		}
	}
}

func TestIntegrationQueryParams(t *testing.T) {
	app := NewApp()

	app.Get("/search", func(c *Context) {
		query := c.Query("q")
		c.JSON(StatusOK, map[string]string{"query": query})
	})

	ctx := newTestCtx(MethodGet, "/search?q=test")
	app.serveRequest(ctx)

	if ctx.Response.StatusCode() != StatusOK {
		t.Errorf("Expected status %d, got %d", StatusOK, ctx.Response.StatusCode())
	}
}

func TestIntegrationPathParams(t *testing.T) {
	app := NewApp()

	app.Get("/users/:id/posts/:postId", func(c *Context) {
		id := c.Param("id")
		postId := c.Param("postId")
		c.JSON(StatusOK, map[string]string{"userId": id, "postId": postId})
	})

	ctx := newTestCtx(MethodGet, "/users/123/posts/456")
	app.serveRequest(ctx)

	if ctx.Response.StatusCode() != StatusOK {
		t.Errorf("Expected status %d, got %d", StatusOK, ctx.Response.StatusCode())
	}
}

func TestIntegrationCookie(t *testing.T) {
	app := NewApp()

	app.Get("/set-cookie", func(c *Context) {
		c.SetCookie("session", "abc123")
		c.Text(StatusOK, "cookie set")
	})

	ctx := newTestCtx(MethodGet, "/set-cookie")
	app.serveRequest(ctx)

	cookie := string(ctx.Response.Header.Peek("Set-Cookie"))
	if cookie == "" {
		t.Error("Expected Set-Cookie header")
	}
}

func TestIntegrationRedirect(t *testing.T) {
	app := NewApp()

	app.Get("/old", func(c *Context) {
		c.Redirect(StatusMovedPermanently, "/new")
	})

	ctx := newTestCtx(MethodGet, "/old")
	app.serveRequest(ctx)

	if ctx.Response.StatusCode() != StatusMovedPermanently {
		t.Errorf("Expected status %d, got %d", StatusMovedPermanently, ctx.Response.StatusCode())
	}
	location := string(ctx.Response.Header.Peek("Location"))
	if location == "" {
		t.Error("Expected Location header to be set")
	}
}

func TestIntegrationHeaders(t *testing.T) {
	app := NewApp()

	app.Get("/headers", func(c *Context) {
		c.SetHeader("X-Custom", "value")
		c.JSON(StatusOK, map[string]string{"header": "set"})
	})

	ctx := newTestCtx(MethodGet, "/headers")
	app.serveRequest(ctx)

	header := string(ctx.Response.Header.Peek("X-Custom"))
	if header != "value" {
		t.Errorf("Expected X-Custom header 'value', got '%s'", header)
	}
}

func TestIntegrationContentNegotiation(t *testing.T) {
	app := NewApp()

	app.Get("/text", func(c *Context) {
		c.Text(StatusOK, "plain text")
	})

	app.Get("/json", func(c *Context) {
		c.JSON(StatusOK, map[string]string{"key": "value"})
	})

	ctx := newTestCtx(MethodGet, "/text")
	app.serveRequest(ctx)
	if ctx.Response.StatusCode() != StatusOK {
		t.Errorf("Expected status %d for text, got %d", StatusOK, ctx.Response.StatusCode())
	}

	ctx = newTestCtx(MethodGet, "/json")
	app.serveRequest(ctx)
	if ctx.Response.StatusCode() != StatusOK {
		t.Errorf("Expected status %d for json, got %d", StatusOK, ctx.Response.StatusCode())
	}
}

func TestIntegrationErrorHandling(t *testing.T) {
	app := NewApp()

	app.Get("/error", func(c *Context) {
		c.JSONError(StatusBadRequest, "invalid request")
	})

	ctx := newTestCtx(MethodGet, "/error")
	app.serveRequest(ctx)

	if ctx.Response.StatusCode() != StatusBadRequest {
		t.Errorf("Expected status %d, got %d", StatusBadRequest, ctx.Response.StatusCode())
	}
}

func TestIntegrationSuccessFail(t *testing.T) {
	app := NewApp()

	app.Get("/success", func(c *Context) {
		c.Success(map[string]string{"result": "ok"})
	})

	app.Get("/fail", func(c *Context) {
		c.Fail(1001, "operation failed")
	})

	ctx := newTestCtx(MethodGet, "/success")
	app.serveRequest(ctx)
	if ctx.Response.StatusCode() != StatusOK {
		t.Errorf("Expected status %d for success, got %d", StatusOK, ctx.Response.StatusCode())
	}

	ctx = newTestCtx(MethodGet, "/fail")
	app.serveRequest(ctx)
	if ctx.Response.StatusCode() != StatusOK {
		t.Errorf("Expected status %d for fail, got %d", StatusOK, ctx.Response.StatusCode())
	}
}

func TestIntegrationNestedGroups(t *testing.T) {
	app := NewApp()

	api := app.Group("/api")
	v1 := api.Group("/v1")
	v2 := api.Group("/v2")

	v1.Get("/resource", func(c *Context) {
		c.Text(StatusOK, "v1 resource")
	})

	v2.Get("/resource", func(c *Context) {
		c.Text(StatusOK, "v2 resource")
	})

	ctx := newTestCtx(MethodGet, "/api/v1/resource")
	app.serveRequest(ctx)
	if ctx.Response.StatusCode() != StatusOK {
		t.Errorf("Expected status %d for v1, got %d", StatusOK, ctx.Response.StatusCode())
	}

	ctx = newTestCtx(MethodGet, "/api/v2/resource")
	app.serveRequest(ctx)
	if ctx.Response.StatusCode() != StatusOK {
		t.Errorf("Expected status %d for v2, got %d", StatusOK, ctx.Response.StatusCode())
	}
}

func TestIntegrationGroupMiddleware(t *testing.T) {
	app := NewApp()
	middlewareOrder := []string{}

	api := app.Group("/api")
	api.Use(func(c *Context) {
		middlewareOrder = append(middlewareOrder, "api-middleware")
		c.Next()
	})

	users := api.Group("/users")
	users.Use(func(c *Context) {
		middlewareOrder = append(middlewareOrder, "users-middleware")
		c.Next()
	})

	users.Get("/:id", func(c *Context) {
		middlewareOrder = append(middlewareOrder, "users-handler")
	})

	ctx := newTestCtx(MethodGet, "/api/users/123")
	app.serveRequest(ctx)

	expected := []string{"api-middleware", "users-middleware", "users-handler"}
	if len(middlewareOrder) != len(expected) {
		t.Errorf("Expected %d middleware calls, got %d: %v", len(expected), len(middlewareOrder), middlewareOrder)
	}
	for i, e := range expected {
		if middlewareOrder[i] != e {
			t.Errorf("Call[%d] = %s, want %s", i, middlewareOrder[i], e)
		}
	}
}
