package lightning

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"text/template"

	"github.com/valyala/fasthttp"
)

func createFasthttpRequest(method, path string) *fasthttp.RequestCtx {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod(method)
	ctx.Request.Header.SetRequestURI(path)
	return ctx
}

func TestNewApp(t *testing.T) {
	app := NewApp()
	if app == nil {
		t.Error("NewApp returned nil")
	}
}

func TestDefaultApp(t *testing.T) {
	app := DefaultApp()

	if app.Logger == nil {
		t.Errorf("Expected Logger field to not be nil")
	}

	if len(app.middlewares) != 2 {
		t.Errorf("Expected 2 middleware functions, but got %d", len(app.middlewares))
	}
}

func TestUse(t *testing.T) {
	app := NewApp()

	mw1 := func(c *Context) {}
	mw2 := func(c *Context) {}

	app.Use(mw1, mw2)

	if len(app.middlewares) != 2 {
		t.Errorf("Expected 2 middleware functions, but got %d", len(app.middlewares))
	}
}

func TestAddRoute(t *testing.T) {
	app := NewApp()
	app.AddRoute(MethodGet, "/test", []HandlerFunc{})
	route, _ := app.router.findRoute(MethodGet, "/test")
	if route == nil {
		t.Errorf("Expected route to be added to router")
	}
}

func TestGetRoute(t *testing.T) {
	app := NewApp()
	app.Get("/test", func(c *Context) {
		c.Text(StatusOK, "Hello, World!")
	})

	ctx := createFasthttpRequest(MethodGet, "/test")
	app.serveRequest(ctx)

	if ctx.Response.StatusCode() != StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			ctx.Response.StatusCode(), StatusOK)
	}

	expected := "Hello, World!"
	if string(ctx.Response.Body()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			string(ctx.Response.Body()), expected)
	}
}

func TestPostRoute(t *testing.T) {
	app := NewApp()
	app.Post("/test", func(c *Context) {
		c.Text(StatusOK, "Hello, World!")
	})

	ctx := createFasthttpRequest(MethodPost, "/test")
	app.serveRequest(ctx)

	if ctx.Response.StatusCode() != StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			ctx.Response.StatusCode(), StatusOK)
	}
}

func TestPutRoute(t *testing.T) {
	app := NewApp()
	app.Put("/test", func(c *Context) {
		c.Text(StatusOK, "Hello, World!")
	})

	ctx := createFasthttpRequest(MethodPut, "/test")
	app.serveRequest(ctx)

	if ctx.Response.StatusCode() != StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			ctx.Response.StatusCode(), StatusOK)
	}
}

func TestDeleteRoute(t *testing.T) {
	app := NewApp()
	app.Delete("/test", func(c *Context) {
		c.Text(StatusOK, "Hello, World!")
	})

	ctx := createFasthttpRequest(MethodDelete, "/test")
	app.serveRequest(ctx)

	if ctx.Response.StatusCode() != StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			ctx.Response.StatusCode(), StatusOK)
	}
}

func TestHeadRoute(t *testing.T) {
	app := NewApp()
	app.Head("/test", func(c *Context) {
		c.Text(StatusOK, "")
	})

	ctx := createFasthttpRequest(MethodHead, "/test")
	app.serveRequest(ctx)

	if ctx.Response.StatusCode() != StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			ctx.Response.StatusCode(), StatusOK)
	}
}

func TestPatchRoute(t *testing.T) {
	app := NewApp()
	app.Patch("/test", func(c *Context) {
		c.Text(StatusOK, "Hello, World!")
	})

	ctx := createFasthttpRequest(MethodPatch, "/test")
	app.serveRequest(ctx)

	if ctx.Response.StatusCode() != StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			ctx.Response.StatusCode(), StatusOK)
	}
}

func TestOptionsRoute(t *testing.T) {
	app := NewApp()
	app.Options("/test", func(c *Context) {
		c.Text(StatusOK, "Hello, World!")
	})

	ctx := createFasthttpRequest(MethodOptions, "/test")
	app.serveRequest(ctx)

	if ctx.Response.StatusCode() != StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			ctx.Response.StatusCode(), StatusOK)
	}
}

func TestNotFoundHandler(t *testing.T) {
	app := NewApp()
	app.Get("/exists", func(c *Context) {
		c.Text(StatusOK, "exists")
	})

	ctx := createFasthttpRequest(MethodGet, "/notfound")
	app.serveRequest(ctx)

	if ctx.Response.StatusCode() != StatusNotFound {
		t.Errorf("Expected status %d, got %d", StatusNotFound, ctx.Response.StatusCode())
	}
}

func TestMiddlewareExecution(t *testing.T) {
	app := NewApp()
	order := []int{}

	app.Use(func(c *Context) {
		order = append(order, 1)
		c.Next()
		order = append(order, 4)
	})

	app.Use(func(c *Context) {
		order = append(order, 2)
		c.Next()
		order = append(order, 3)
	})

	app.Get("/test", func(c *Context) {
		order = append(order, 5)
	})

	ctx := createFasthttpRequest(MethodGet, "/test")
	app.serveRequest(ctx)

	expected := []int{1, 2, 5, 3, 4}
	if !stringslicesEqual(order, expected) {
		t.Errorf("Middleware execution order wrong: got %v want %v", order, expected)
	}
}

func stringslicesEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestStaticFiles(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "static_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	tmpFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(tmpFile, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	app := NewApp()
	app.Static(tmpDir, "/static")

	ctx := createFasthttpRequest(MethodGet, "/static/test.txt")
	app.serveRequest(ctx)

	if ctx.Response.StatusCode() != StatusOK {
		t.Errorf("Expected status %d, got %d", StatusOK, ctx.Response.StatusCode())
	}
}

func TestStaticFilesNotFound(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "static_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	app := NewApp()
	app.Static(tmpDir, "/static")

	ctx := createFasthttpRequest(MethodGet, "/static/nonexistent.txt")
	app.serveRequest(ctx)

	if ctx.Response.StatusCode() != StatusNotFound {
		t.Errorf("Expected status %d, got %d", StatusNotFound, ctx.Response.StatusCode())
	}
}

func TestGroupRoute(t *testing.T) {
	app := NewApp()
	group := app.Group("/api")

	group.Get("/test", func(c *Context) {
		c.Text(StatusOK, "group route")
	})

	ctx := createFasthttpRequest(MethodGet, "/api/test")
	app.serveRequest(ctx)

	if ctx.Response.StatusCode() != StatusOK {
		t.Errorf("Expected status %d, got %d", StatusOK, ctx.Response.StatusCode())
	}
}

func TestGroupMiddleware(t *testing.T) {
	app := NewApp()
	order := []int{}

	group := app.Group("/api")
	group.Use(func(c *Context) {
		order = append(order, 1)
		c.Next()
	})

	group.Get("/test", func(c *Context) {
		order = append(order, 2)
	})

	ctx := createFasthttpRequest(MethodGet, "/api/test")
	app.serveRequest(ctx)

	expected := []int{1, 2}
	if !stringslicesEqual(order, expected) {
		t.Errorf("Middleware execution order wrong: got %v want %v", order, expected)
	}
}

func TestNestedGroup(t *testing.T) {
	app := NewApp()
	group := app.Group("/api")
	nested := group.Group("/v1")

	nested.Get("/test", func(c *Context) {
		c.Text(StatusOK, "nested group route")
	})

	ctx := createFasthttpRequest(MethodGet, "/api/v1/test")
	app.serveRequest(ctx)

	if ctx.Response.StatusCode() != StatusOK {
		t.Errorf("Expected status %d, got %d", StatusOK, ctx.Response.StatusCode())
	}
}

func TestConfigMerge(t *testing.T) {
	config1 := &Config{
		AppName: "app1",
	}
	config2 := &Config{
		AppName: "app2",
	}

	merged := config1.merge(config2)
	if merged.AppName != "app2" {
		t.Errorf("Expected AppName 'app2', got '%s'", merged.AppName)
	}
}

func TestLoadHTMLGlob(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "templates_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	tmplPath := filepath.Join(tmpDir, "test.html")
	if err := os.WriteFile(tmplPath, []byte("<html>{{.Name}}</html>"), 0644); err != nil {
		t.Fatal(err)
	}

	app := NewApp()
	app.SetFuncMap(template.FuncMap{})
	app.LoadHTMLGlob(filepath.Join(tmpDir, "*.html"))

	if app.htmlTemplates == nil {
		t.Error("Expected htmlTemplates to be set")
	}
}

func TestJSONResponse(t *testing.T) {
	app := NewApp()
	app.Get("/test", func(c *Context) {
		c.JSON(StatusOK, map[string]string{"message": "hello"})
	})

	ctx := createFasthttpRequest(MethodGet, "/test")
	app.serveRequest(ctx)

	if ctx.Response.StatusCode() != StatusOK {
		t.Errorf("Expected status %d, got %d", StatusOK, ctx.Response.StatusCode())
	}

	contentType := string(ctx.Response.Header.ContentType())
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("Expected Content-Type to contain 'application/json', got '%s'", contentType)
	}
}

func TestRedirect(t *testing.T) {
	app := NewApp()
	app.Get("/test", func(c *Context) {
		c.Redirect(StatusMovedPermanently, "/new")
	})

	ctx := createFasthttpRequest(MethodGet, "/test")
	app.serveRequest(ctx)

	if ctx.Response.StatusCode() != StatusMovedPermanently {
		t.Errorf("Expected status %d, got %d", StatusMovedPermanently, ctx.Response.StatusCode())
	}

	location := string(ctx.Response.Header.Peek("Location"))
	if location == "" {
		t.Error("Expected Location header to be set")
	}
}

func TestRequestHandler(t *testing.T) {
	app := NewApp()
	app.Get("/test", func(c *Context) {
		c.Text(StatusOK, "handler")
	})

	handler := app.RequestHandler()
	if handler == nil {
		t.Error("RequestHandler returned nil")
	}

	ctx := createFasthttpRequest(MethodGet, "/test")
	handler(ctx)

	if ctx.Response.StatusCode() != StatusOK {
		t.Errorf("Expected status %d, got %d", StatusOK, ctx.Response.StatusCode())
	}
}

func TestAcquireReleaseContext(t *testing.T) {
	app := NewApp()

	ctx := createFasthttpRequest(MethodGet, "/test")
	c := app.acquireContext(ctx)

	if c == nil {
		t.Fatal("acquireContext returned nil")
	}

	if c.ctx != ctx {
		t.Error("Context not set correctly")
	}

	app.releaseContext(c)
}
