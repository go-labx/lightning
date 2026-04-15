package lightning

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"text/template"
	"time"

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

	if len(app.middlewares) != 3 {
		t.Errorf("Expected 3 middleware functions, but got %d", len(app.middlewares))
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

func TestConfigMergeJSONEncoder(t *testing.T) {
	cfg := &Config{}
	encoder := func(v interface{}) ([]byte, error) { return []byte("{}"), nil }
	merged := cfg.merge(&Config{JSONEncoder: encoder})
	if merged.JSONEncoder == nil {
		t.Error("Expected JSONEncoder to be set")
	}
}

func TestConfigMergeJSONDecoder(t *testing.T) {
	cfg := &Config{}
	decoder := func(data []byte, v interface{}) error { return nil }
	merged := cfg.merge(&Config{JSONDecoder: decoder})
	if merged.JSONDecoder == nil {
		t.Error("Expected JSONDecoder to be set")
	}
}

func TestConfigMergeNotFoundHandler(t *testing.T) {
	cfg := &Config{}
	handler := func(c *Context) {}
	merged := cfg.merge(&Config{NotFoundHandler: handler})
	if merged.NotFoundHandler == nil {
		t.Error("Expected NotFoundHandler to be set")
	}
}

func TestConfigMergeMaxRequestBodySize(t *testing.T) {
	cfg := &Config{}
	merged := cfg.merge(&Config{MaxRequestBodySize: 4096})
	if merged.MaxRequestBodySize != 4096 {
		t.Errorf("Expected MaxRequestBodySize 4096, got %d", merged.MaxRequestBodySize)
	}
}

func TestConfigMergeMaxRequestBodySizeZero(t *testing.T) {
	cfg := &Config{MaxRequestBodySize: 1024}
	merged := cfg.merge(&Config{MaxRequestBodySize: 0})
	if merged.MaxRequestBodySize != 1024 {
		t.Errorf("Expected MaxRequestBodySize 1024, got %d", merged.MaxRequestBodySize)
	}
}

func TestConfigMergeNilInMiddle(t *testing.T) {
	cfg := &Config{AppName: "original"}
	merged := cfg.merge(&Config{AppName: "first"}, nil, &Config{AppName: "second"})
	if merged.AppName != "second" {
		t.Errorf("Expected AppName 'second', got '%s'", merged.AppName)
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

func TestConfigMergeAll(t *testing.T) {
	cfg := &Config{}

	merged := cfg.merge(&Config{
		AppName:            "test-app",
		EnableDebug:        true,
		MaxRequestBodySize: 1024,
	})

	if merged.AppName != "test-app" {
		t.Errorf("Expected AppName 'test-app', got '%s'", merged.AppName)
	}
	if !merged.EnableDebug {
		t.Error("Expected EnableDebug to be true")
	}
	if merged.MaxRequestBodySize != 1024 {
		t.Errorf("Expected MaxRequestBodySize 1024, got %d", merged.MaxRequestBodySize)
	}
}

func TestConfigMergeNil(t *testing.T) {
	cfg := &Config{AppName: "original"}

	merged := cfg.merge(nil, nil)

	if merged.AppName != "original" {
		t.Errorf("Expected AppName 'original', got '%s'", merged.AppName)
	}
}

func TestConfigMergePartial(t *testing.T) {
	cfg := &Config{AppName: "original", EnableDebug: false}

	merged := cfg.merge(&Config{AppName: ""}, &Config{EnableDebug: true})

	if merged.AppName != "original" {
		t.Errorf("Expected AppName 'original', got '%s'", merged.AppName)
	}
	if !merged.EnableDebug {
		t.Error("Expected EnableDebug to be true")
	}
}

func TestNewAppWithConfig(t *testing.T) {
	app := NewApp(&Config{
		AppName:            "custom-app",
		EnableDebug:        true,
		MaxRequestBodySize: 2048,
	})

	if app.Config.AppName != "custom-app" {
		t.Errorf("Expected AppName 'custom-app', got '%s'", app.Config.AppName)
	}
	if !app.Config.EnableDebug {
		t.Error("Expected EnableDebug to be true")
	}
	if app.Config.MaxRequestBodySize != 2048 {
		t.Errorf("Expected MaxRequestBodySize 2048, got %d", app.Config.MaxRequestBodySize)
	}
}

func TestNewAppWithMultipleConfigs(t *testing.T) {
	app := NewApp(
		&Config{AppName: "first"},
		&Config{AppName: "second"},
	)

	if app.Config.AppName != "second" {
		t.Errorf("Expected AppName 'second', got '%s'", app.Config.AppName)
	}
}

func TestNewAppDebugRoute(t *testing.T) {
	app := NewApp(&Config{EnableDebug: true})

	ctx := createFasthttpRequest(MethodGet, "/__debug__/router_map")
	app.serveRequest(ctx)

	if ctx.Response.StatusCode() != StatusOK {
		t.Errorf("Expected status %d, got %d", StatusOK, ctx.Response.StatusCode())
	}
}

func TestDefaultJSONUnmarshal(t *testing.T) {
	data := []byte(`{"key":"value"}`)
	var result map[string]string

	err := defaultJSONUnmarshal(data, &result)
	if err != nil {
		t.Fatalf("defaultJSONUnmarshal returned error: %v", err)
	}
	if result["key"] != "value" {
		t.Errorf("Expected 'value', got '%s'", result["key"])
	}
}

func TestDefaultJSONMarshal(t *testing.T) {
	data := map[string]string{"key": "value"}

	result, err := defaultJSONMarshal(data)
	if err != nil {
		t.Fatalf("defaultJSONMarshal returned error: %v", err)
	}
	if string(result) != `{"key":"value"}` {
		t.Errorf("Expected '{\"key\":\"value\"}', got '%s'", string(result))
	}
}

func TestStaticAbsolute(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "static_abs")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	if err := os.WriteFile(tmpDir+"/file.txt", []byte("content"), 0644); err != nil {
		t.Fatal(err)
	}

	app := NewApp()
	app.Static(tmpDir, "/static")

	ctx := createFasthttpRequest(MethodGet, "/static/file.txt")
	app.serveRequest(ctx)

	if ctx.Response.StatusCode() != StatusOK {
		t.Errorf("Expected status %d, got %d", StatusOK, ctx.Response.StatusCode())
	}
}

func TestStaticNotFound(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "static_notfound")
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

func TestRouterSearchNotFound(t *testing.T) {
	router := newRouter()
	router.addRoute(MethodGet, "/api/users", []HandlerFunc{})

	route, params := router.findRoute(MethodGet, "/api/posts")
	if route != nil {
		t.Error("Expected nil route for non-matching path")
	}
	if params != nil {
		t.Error("Expected nil params for non-matching path")
	}
}

func TestRouterSearchWildcard(t *testing.T) {
	router := newRouter()
	router.addRoute(MethodGet, "/api/*filepath", []HandlerFunc{func(c *Context) {}})

	route, params := router.findRoute(MethodGet, "/api/users/123/posts")
	if route == nil {
		t.Fatal("Expected route for wildcard path")
	}
	if params["filepath"] != "users/123/posts" {
		t.Errorf("Expected wildcard param 'users/123/posts', got '%s'", params["filepath"])
	}
}

func TestRouterSearchParam(t *testing.T) {
	router := newRouter()
	router.addRoute(MethodGet, "/users/:id", []HandlerFunc{func(c *Context) {}})

	route, params := router.findRoute(MethodGet, "/users/42")
	if route == nil {
		t.Fatal("Expected route for param path")
	}
	if params["id"] != "42" {
		t.Errorf("Expected param id=42, got '%s'", params["id"])
	}
}

func TestRouterFindRouteNotFound(t *testing.T) {
	router := newRouter()
	router.addRoute(MethodGet, "/test", []HandlerFunc{})

	route, params := router.findRoute(MethodPost, "/test")
	if route != nil {
		t.Error("Expected nil route for wrong method")
	}
	if params != nil {
		t.Error("Expected nil params for wrong method")
	}
}

func TestRouterFindRouteEmpty(t *testing.T) {
	router := newRouter()

	route, params := router.findRoute(MethodGet, "/nonexistent")
	if route != nil {
		t.Error("Expected nil route for empty router")
	}
	if params != nil {
		t.Error("Expected nil params for empty router")
	}
}

func TestRouterSearchMultipleParams(t *testing.T) {
	router := newRouter()
	router.addRoute(MethodGet, "/users/:userId/posts/:postId", []HandlerFunc{func(c *Context) {}})

	route, params := router.findRoute(MethodGet, "/users/1/posts/2")
	if route == nil {
		t.Fatal("Expected route for multi-param path")
	}
	if params["userId"] != "1" {
		t.Errorf("Expected userId=1, got '%s'", params["userId"])
	}
	if params["postId"] != "2" {
		t.Errorf("Expected postId=2, got '%s'", params["postId"])
	}
}

func TestLogger(t *testing.T) {
	app := NewApp()
	if app.Logger == nil {
		t.Error("Expected Logger to be non-nil")
	}
}

func TestRouterSearchEmptyPath(t *testing.T) {
	router := newRouter()
	router.addRoute(MethodGet, "/", []HandlerFunc{func(c *Context) {}})

	route, _ := router.findRoute(MethodGet, "/")
	if route == nil {
		t.Error("Expected route for root path")
	}
}

func TestRouterSearchMethodNotAllowed(t *testing.T) {
	router := newRouter()
	router.addRoute(MethodGet, "/test", []HandlerFunc{func(c *Context) {}})

	route, _ := router.findRoute(MethodPost, "/test")
	if route != nil {
		t.Error("Expected nil route for method not allowed")
	}
}

func TestContext_NextOutOfBounds(t *testing.T) {
	c, _ := createTestContext("GET", "/test", nil)
	c.handlers = []HandlerFunc{}
	c.index = 0

	c.Next()

	if c.index != 1 {
		t.Errorf("Expected index 1, got %d", c.index)
	}
}

func TestRouterSearchParamWithMultipleLevels(t *testing.T) {
	router := newRouter()
	router.addRoute(MethodGet, "/api/:version/users/:id", []HandlerFunc{func(c *Context) {}})

	route, params := router.findRoute(MethodGet, "/api/v1/users/42")
	if route == nil {
		t.Fatal("Expected route for multi-level param path")
	}
	if params["version"] != "v1" {
		t.Errorf("Expected version=v1, got '%s'", params["version"])
	}
	if params["id"] != "42" {
		t.Errorf("Expected id=42, got '%s'", params["id"])
	}
}

func TestRouterAddRouteMultiple(t *testing.T) {
	router := newRouter()
	router.addRoute(MethodGet, "/a", []HandlerFunc{func(c *Context) {}})
	router.addRoute(MethodPost, "/a", []HandlerFunc{func(c *Context) {}})
	router.addRoute(MethodPut, "/a", []HandlerFunc{func(c *Context) {}})

	if len(router.Roots) != 3 {
		t.Errorf("Expected 3 method roots, got %d", len(router.Roots))
	}
}

func TestRouterSearchNotFoundMethod(t *testing.T) {
	router := newRouter()
	router.addRoute(MethodGet, "/test", []HandlerFunc{func(c *Context) {}})

	route, params := router.findRoute(MethodDelete, "/test")
	if route != nil {
		t.Error("Expected nil route for non-existent method")
	}
	if params != nil {
		t.Error("Expected nil params for non-existent method")
	}
}

func TestParsePatternWithWildcard(t *testing.T) {
	parts := parsePattern("/files/*filepath")
	if len(parts) != 2 {
		t.Errorf("Expected 2 parts, got %d", len(parts))
	}
	if parts[1] != "*filepath" {
		t.Errorf("Expected '*filepath', got '%s'", parts[1])
	}
}

func TestResolveAddressWithPortEnv(t *testing.T) {
	os.Setenv("PORT", "8080")
	defer os.Unsetenv("PORT")

	addr := resolveAddress([]string{})
	if addr != ":8080" {
		t.Errorf("Expected ':8080', got '%s'", addr)
	}
}

func TestResolveAddressSingleParam(t *testing.T) {
	os.Unsetenv("PORT")
	addr := resolveAddress([]string{":9090"})
	if addr != ":9090" {
		t.Errorf("Expected ':9090', got '%s'", addr)
	}
}

func TestResolveAddressMultipleParam(t *testing.T) {
	os.Unsetenv("PORT")
	defer func() {
		recover()
	}()
	resolveAddress([]string{":8080", ":9090"})
}

func TestDefaultNotFoundHandler(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod("GET")
	ctx.Request.Header.SetRequestURI("/notfound")

	c := &Context{
		ctx:   ctx,
		index: -1,
		res:   newResponse(ctx),
	}
	defaultNotFound(c)
	c.flush()

	if ctx.Response.StatusCode() != StatusNotFound {
		t.Errorf("Expected status %d, got %d", StatusNotFound, ctx.Response.StatusCode())
	}
}

func TestDefaultInternalServerErrorHandler(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod("GET")
	ctx.Request.Header.SetRequestURI("/error")

	c := &Context{
		ctx:   ctx,
		index: -1,
		res:   newResponse(ctx),
	}
	defaultInternalServerError(c)
	c.flush()

	if ctx.Response.StatusCode() != StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", StatusInternalServerError, ctx.Response.StatusCode())
	}
}

func TestNewRouter(t *testing.T) {
	router := newRouter()
	if router.Roots == nil {
		t.Error("Expected Roots to be initialized")
	}
}

func TestNodeMatchChild(t *testing.T) {
	n := &node{
		Children: map[string]*node{
			"foo": {Part: "foo"},
		},
	}

	child := n.matchChild("foo")
	if child == nil {
		t.Error("Expected child node")
	}

	child = n.matchChild("bar")
	if child != nil {
		t.Error("Expected nil for non-existent child")
	}
}

func TestNodeInsert(t *testing.T) {
	root := &node{}
	root.insert("/a/b/c", []string{"a", "b", "c"}, 0, []HandlerFunc{func(c *Context) {}})

	if root.Children == nil {
		t.Fatal("Expected children to be initialized")
	}

	n := root.matchChild("a")
	if n == nil {
		t.Fatal("Expected node 'a'")
	}
	if n.Part != "a" {
		t.Errorf("Expected part 'a', got '%s'", n.Part)
	}
}

func TestNodeInsertWildParam(t *testing.T) {
	root := &node{}
	root.insert("/users/:id", []string{"users", ":id"}, 0, []HandlerFunc{func(c *Context) {}})

	users := root.matchChild("users")
	if users == nil {
		t.Fatal("Expected 'users' node")
	}

	idNode := users.matchChild(":id")
	if idNode == nil {
		t.Fatal("Expected ':id' node")
	}
	if !idNode.IsWild {
		t.Error("Expected :id to be wild")
	}
}

func TestNodeSearch(t *testing.T) {
	root := &node{}
	root.insert("/a/b", []string{"a", "b"}, 0, []HandlerFunc{func(c *Context) {}})

	n := root.search([]string{"a", "b"}, 0)
	if n == nil {
		t.Error("Expected to find node")
	}
	if n.Pattern != "/a/b" {
		t.Errorf("Expected pattern '/a/b', got '%s'", n.Pattern)
	}
}

func TestNodeSearchNotFound(t *testing.T) {
	root := &node{}
	root.insert("/a/b", []string{"a", "b"}, 0, []HandlerFunc{func(c *Context) {}})

	n := root.search([]string{"a", "c"}, 0)
	if n != nil {
		t.Error("Expected nil for non-existent path")
	}
}

func TestNodeSearchWild(t *testing.T) {
	root := &node{}
	root.insert("/api/*", []string{"api", "*"}, 0, []HandlerFunc{func(c *Context) {}})

	n := root.search([]string{"api", "users", "123"}, 0)
	if n == nil {
		t.Error("Expected to find wildcard node")
	}
}

func TestNodeSearchParam(t *testing.T) {
	root := &node{}
	root.insert("/users/:id", []string{"users", ":id"}, 0, []HandlerFunc{func(c *Context) {}})

	n := root.search([]string{"users", "42"}, 0)
	if n == nil {
		t.Error("Expected to find param node")
	}
}

func TestNodeSearchEmptyParts(t *testing.T) {
	root := &node{Pattern: "/"}
	n := root.search([]string{}, 0)
	if n == nil {
		t.Error("Expected to find root node")
	}
}

func TestNodeSearchWildReturnNil(t *testing.T) {
	root := &node{}
	root.insert("/api/*", []string{"api", "*"}, 0, []HandlerFunc{func(c *Context) {}})

	n := root.search([]string{"api"}, 0)
	if n != nil {
		t.Error("Expected nil for incomplete path")
	}
}

func TestResponseFlushWithRedirect(t *testing.T) {
	resp, ctx := createResponse()

	resp.redirect(StatusFound, "/new")
	resp.flush()

	if ctx.Response.StatusCode() != StatusFound {
		t.Errorf("Expected status %d, got %d", StatusFound, ctx.Response.StatusCode())
	}
	location := string(ctx.Response.Header.Peek("Location"))
	if location == "" {
		t.Error("Expected Location header to be set")
	}
}

func TestResponseFlushWithFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "response_test*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte("file content")); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	resp, ctx := createResponse()
	resp.file(tmpFile.Name())
	resp.flush()

	if ctx.Response.StatusCode() != StatusOK {
		t.Errorf("Expected status %d, got %d", StatusOK, ctx.Response.StatusCode())
	}
}

func TestResponseFlushWithCookie(t *testing.T) {
	resp, ctx := createResponse()

	resp.cookies.set("session", "abc")
	resp.flush()

	cookie := string(ctx.Response.Header.Peek("Set-Cookie"))
	if !strings.Contains(cookie, "session=abc") {
		t.Errorf("Expected cookie to contain 'session=abc', got %s", cookie)
	}
}

func TestResponseFlushWithBody(t *testing.T) {
	resp, ctx := createResponse()

	resp.setBody([]byte("hello"))
	resp.flush()

	if string(ctx.Response.Body()) != "hello" {
		t.Errorf("Expected body 'hello', got '%s'", string(ctx.Response.Body()))
	}
}

func TestResponseFlushWithCookieAndBody(t *testing.T) {
	resp, ctx := createResponse()

	resp.cookies.set("test", "value")
	resp.setBody([]byte("body"))
	resp.flush()

	if string(ctx.Response.Body()) != "body" {
		t.Errorf("Expected body 'body', got '%s'", string(ctx.Response.Body()))
	}
}

func TestResponseFlushWithCookieAndRedirect(t *testing.T) {
	resp, ctx := createResponse()

	resp.cookies.set("session", "abc")
	resp.redirect(StatusFound, "/new")
	resp.flush()

	if ctx.Response.StatusCode() != StatusFound {
		t.Errorf("Expected status %d, got %d", StatusFound, ctx.Response.StatusCode())
	}
}

func TestResponseFileError(t *testing.T) {
	resp, _ := createResponse()

	err := resp.file("/nonexistent/file.txt")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestResponseSetStatus(t *testing.T) {
	resp, _ := createResponse()

	resp.setStatus(StatusCreated)
	if resp.statusCode != StatusCreated {
		t.Errorf("Expected status %d, got %d", StatusCreated, resp.statusCode)
	}
}

func TestResponseSetBody(t *testing.T) {
	resp, _ := createResponse()

	resp.setBody([]byte("test"))
	if string(resp.body) != "test" {
		t.Errorf("Expected body 'test', got '%s'", string(resp.body))
	}
}

func TestResponseRedirect(t *testing.T) {
	resp, _ := createResponse()

	resp.redirect(StatusMovedPermanently, "/new")
	if resp.redirectTo != "/new" {
		t.Errorf("Expected redirectTo '/new', got '%s'", resp.redirectTo)
	}
	if resp.statusCode != StatusMovedPermanently {
		t.Errorf("Expected status %d, got %d", StatusMovedPermanently, resp.statusCode)
	}
}

func TestResponseAddHeader(t *testing.T) {
	resp, ctx := createResponse()

	resp.addHeader("X-Custom", "value")
	hdr := string(ctx.Response.Header.Peek("X-Custom"))
	if hdr != "value" {
		t.Errorf("Expected header 'value', got '%s'", hdr)
	}
}

func TestResponseSetHeader(t *testing.T) {
	resp, ctx := createResponse()

	resp.setHeader("Content-Type", "text/plain")
	hdr := string(ctx.Response.Header.Peek("Content-Type"))
	if hdr != "text/plain" {
		t.Errorf("Expected header 'text/plain', got '%s'", hdr)
	}
}

func TestResponseDelHeader(t *testing.T) {
	resp, ctx := createResponse()

	resp.setHeader("X-Custom", "value")
	resp.delHeader("X-Custom")
	hdr := string(ctx.Response.Header.Peek("X-Custom"))
	if hdr != "" {
		t.Errorf("Expected header to be deleted, got '%s'", hdr)
	}
}

func TestResponseNew(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	resp := newResponse(ctx)

	if resp.ctx != ctx {
		t.Error("Expected ctx to be set")
	}
	if resp.statusCode != StatusNotFound {
		t.Errorf("Expected default status %d, got %d", StatusNotFound, resp.statusCode)
	}
	if resp.cookies == nil {
		t.Error("Expected cookies to be initialized")
	}
}

func TestResponseFlushEmpty(t *testing.T) {
	resp, ctx := createResponse()

	resp.flush()

	if ctx.Response.StatusCode() != StatusNotFound {
		t.Errorf("Expected default status %d, got %d", StatusNotFound, ctx.Response.StatusCode())
	}
}

func TestStaticRelativePath(t *testing.T) {
	app := NewApp()
	app.Static("assets", "/static")

	ctx := createFasthttpRequest(MethodGet, "/static/test.txt")
	app.serveRequest(ctx)

	if ctx.Response.StatusCode() != StatusNotFound {
		t.Errorf("Expected status %d, got %d", StatusNotFound, ctx.Response.StatusCode())
	}
}

func TestShutdownWithNilServer(t *testing.T) {
	app := NewApp()

	app.Shutdown()
}

func TestShutdownWithServer(t *testing.T) {
	app := NewApp()
	app.Get("/test", func(c *Context) {
		c.Text(StatusOK, "ok")
	})

	go func() {
		_ = app.Run(":0")
	}()

	time.Sleep(50 * time.Millisecond)

	app.Shutdown()
}

func TestRunGracefulWithZeroTimeout(t *testing.T) {
	app := NewApp()
	app.Get("/test", func(c *Context) {
		c.Text(StatusOK, "ok")
	})

	go func() {
		_ = app.RunGraceful(0, ":0")
	}()

	time.Sleep(50 * time.Millisecond)

	app.Shutdown()
}

func TestContextResetClearsData(t *testing.T) {
	app := NewApp()
	ctx := createFasthttpRequest(MethodGet, "/test")
	c := app.acquireContext(ctx)

	c.SetData("key", "value")
	if c.GetData("key") != "value" {
		t.Error("Expected data to be set")
	}

	app.releaseContext(c)

	c2 := app.acquireContext(createFasthttpRequest(MethodGet, "/test2"))
	if c2.GetData("key") != nil {
		t.Error("Expected data to be cleared after reset")
	}
	app.releaseContext(c2)
}

func TestServeRequestWithNotFoundRoute(t *testing.T) {
	app := NewApp()
	app.Use(func(c *Context) {
		c.Next()
	})

	ctx := createFasthttpRequest(MethodGet, "/nonexistent")
	app.serveRequest(ctx)

	if ctx.Response.StatusCode() != StatusNotFound {
		t.Errorf("Expected status %d, got %d", StatusNotFound, ctx.Response.StatusCode())
	}
}

func TestHelmet(t *testing.T) {
	app := NewApp()
	app.Use(Helmet())
	app.Get("/test", func(c *Context) {
		c.Text(StatusOK, "ok")
	})

	ctx := createFasthttpRequest(MethodGet, "/test")
	app.serveRequest(ctx)

	tests := []struct {
		header string
		want   string
	}{
		{"X-Content-Type-Options", "nosniff"},
		{"X-Frame-Options", "DENY"},
		{"X-XSS-Protection", "1; mode=block"},
		{"Referrer-Policy", "strict-origin-when-cross-origin"},
	}

	for _, tt := range tests {
		got := string(ctx.Response.Header.Peek(tt.header))
		if got != tt.want {
			t.Errorf("Expected %s=%q, got %q", tt.header, tt.want, got)
		}
	}
}

func TestStaticPathTraversalBlocked(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "safe.txt"), []byte("safe"), 0644)

	app := NewApp()
	app.Static(tmpDir, "/static")

	ctx := createFasthttpRequest(MethodGet, "/static/../etc/passwd")
	app.serveRequest(ctx)

	if ctx.Response.StatusCode() != StatusNotFound && ctx.Response.StatusCode() != StatusForbidden {
		t.Errorf("Expected 404 or 403 for path traversal, got %d", ctx.Response.StatusCode())
	}

	body := string(ctx.Response.Body())
	if strings.Contains(body, "root:") {
		t.Error("Path traversal succeeded — /etc/passwd was read!")
	}
}

func TestStaticPathEscapeBlocked(t *testing.T) {
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "assets")
	os.Mkdir(subDir, 0755)

	app := NewApp()
	app.Static(subDir, "/static")

	ctx := createFasthttpRequest(MethodGet, "/static/../../../etc/passwd")
	app.serveRequest(ctx)

	if ctx.Response.StatusCode() != StatusNotFound && ctx.Response.StatusCode() != StatusForbidden {
		t.Errorf("Expected 404 or 403 for path escape, got %d", ctx.Response.StatusCode())
	}

	body := string(ctx.Response.Body())
	if strings.Contains(body, "root:") {
		t.Error("Path escape succeeded — /etc/passwd was read!")
	}
}

func TestStaticDirectoryListingBlocked(t *testing.T) {
	tmpDir := t.TempDir()
	os.Mkdir(filepath.Join(tmpDir, "subdir"), 0755)

	app := NewApp()
	app.Static(tmpDir, "/static")

	ctx := createFasthttpRequest(MethodGet, "/static/subdir")
	app.serveRequest(ctx)

	if ctx.Response.StatusCode() != StatusNotFound {
		t.Errorf("Expected status %d for directory, got %d", StatusNotFound, ctx.Response.StatusCode())
	}
}

func TestParseTrustedProxies_Valid(t *testing.T) {
	app := NewApp(&Config{TrustedProxies: []string{"192.168.1.0/24", "10.0.0.1"}})

	if len(app.trustedProxies) != 2 {
		t.Fatalf("Expected 2 trusted proxies, got %d", len(app.trustedProxies))
	}

	if !app.isTrustedProxy("192.168.1.100:8080") {
		t.Error("Expected 192.168.1.100 to be trusted")
	}
	if app.isTrustedProxy("192.168.2.1:8080") {
		t.Error("Expected 192.168.2.1 to NOT be trusted")
	}
	if !app.isTrustedProxy("10.0.0.1:8080") {
		t.Error("Expected 10.0.0.1 to be trusted")
	}
}

func TestParseTrustedProxies_InvalidCIDR(t *testing.T) {
	app := NewApp(&Config{TrustedProxies: []string{"not-a-valid-cidr"}})

	if len(app.trustedProxies) != 0 {
		t.Errorf("Expected 0 trusted proxies for invalid CIDR, got %d", len(app.trustedProxies))
	}
}

func TestParseTrustedProxies_Nil(t *testing.T) {
	app := NewApp()

	if len(app.trustedProxies) != 0 {
		t.Errorf("Expected 0 trusted proxies for nil config, got %d", len(app.trustedProxies))
	}
}

func TestIsTrustedProxy_NoProxies(t *testing.T) {
	app := NewApp()

	if app.isTrustedProxy("127.0.0.1:8080") {
		t.Error("Expected no proxy to be trusted when TrustedProxies is empty")
	}
}

func TestIsTrustedProxy_InvalidIP(t *testing.T) {
	app := NewApp(&Config{TrustedProxies: []string{"127.0.0.1"}})

	if app.isTrustedProxy("not-an-ip") {
		t.Error("Expected invalid IP to not be trusted")
	}
}

func TestIsTrustedProxy_NoPort(t *testing.T) {
	app := NewApp(&Config{TrustedProxies: []string{"127.0.0.1"}})

	if !app.isTrustedProxy("127.0.0.1") {
		t.Error("Expected 127.0.0.1 without port to be trusted")
	}
}

func TestDebugEndpointWithToken(t *testing.T) {
	app := NewApp(&Config{
		EnableDebug: true,
		DebugToken:  "secret123",
	})

	ctx := createFasthttpRequest(MethodGet, "/__debug__/router_map")
	app.serveRequest(ctx)

	if ctx.Response.StatusCode() != StatusUnauthorized {
		t.Errorf("Expected status %d without token, got %d", StatusUnauthorized, ctx.Response.StatusCode())
	}

	ctx2 := createFasthttpRequest(MethodGet, "/__debug__/router_map?token=secret123")
	app.serveRequest(ctx2)

	if ctx2.Response.StatusCode() != StatusOK {
		t.Errorf("Expected status %d with valid token, got %d", StatusOK, ctx2.Response.StatusCode())
	}

	ctx3 := createFasthttpRequest(MethodGet, "/__debug__/router_map?token=wrong")
	app.serveRequest(ctx3)

	if ctx3.Response.StatusCode() != StatusUnauthorized {
		t.Errorf("Expected status %d with wrong token, got %d", StatusUnauthorized, ctx3.Response.StatusCode())
	}
}

func TestConfigMergeTrustedProxies(t *testing.T) {
	cfg := &Config{}
	merged := cfg.merge(&Config{
		TrustedProxies: []string{"10.0.0.0/8"},
	})

	if len(merged.TrustedProxies) != 1 || merged.TrustedProxies[0] != "10.0.0.0/8" {
		t.Errorf("Expected TrustedProxies to be merged, got %v", merged.TrustedProxies)
	}
}

func TestConfigMergeDebugToken(t *testing.T) {
	cfg := &Config{}
	merged := cfg.merge(&Config{
		DebugToken: "mytoken",
	})

	if merged.DebugToken != "mytoken" {
		t.Errorf("Expected DebugToken 'mytoken', got '%s'", merged.DebugToken)
	}
}

func TestRemoteAddrWithoutTrustedProxies(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod("GET")
	ctx.Request.Header.SetRequestURI("/test")
	ctx.Request.Header.Set("X-Real-IP", "1.2.3.4")

	r := newRequest(ctx)
	r.app = NewApp()
	addr := r.remoteAddr()

	if addr == "1.2.3.4" {
		t.Error("Expected X-Real-IP to be ignored when no TrustedProxies configured")
	}
}
