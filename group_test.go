package lightning

import (
	"reflect"
	"testing"

	"github.com/valyala/fasthttp"
)

func TestNewGroup(t *testing.T) {
	app := NewApp()
	prefix := "/test"
	group := newGroup(app, prefix)

	if group.app != app {
		t.Errorf("Expected app to be %v, but got %v", app, group.app)
	}

	if group.prefix != prefix {
		t.Errorf("Expected prefix to be %v, but got %v", prefix, group.prefix)
	}

	if group.parent != nil {
		t.Errorf("Expected parent to be nil, but got %v", group.parent)
	}

	if len(group.middlewares) != 0 {
		t.Errorf("Expected middlewares to be empty, but got %v", group.middlewares)
	}
}

func TestGetFullPrefix(t *testing.T) {
	app := NewApp()
	group1 := app.Group("/api")
	group2 := group1.Group("/v1")
	group3 := group2.Group("/users")

	expectedPrefix := "/api/v1/users"
	actualPrefix := group3.getFullPrefix()

	if actualPrefix != expectedPrefix {
		t.Errorf("Expected prefix %s, but got %s", expectedPrefix, actualPrefix)
	}
}

func TestGetMiddlewares(t *testing.T) {
	app := NewApp()
	group1 := app.Group("/group1")
	group2 := group1.Group("/group2")
	mw1 := func(c *Context) {}
	mw2 := func(c *Context) {}
	group1.Use(mw1)
	group2.Use(mw2)

	middlewares := group2.getMiddlewares()

	if len(middlewares) != 2 {
		t.Errorf("Expected 2 middlewares, but got %d", len(middlewares))
	}
}

func TestGroup(t *testing.T) {
	app := NewApp()
	group := app.Group("/api")

	if group.prefix != "/api" {
		t.Errorf("Expected prefix to be '/api', but got '%s'", group.prefix)
	}

	if group.parent != nil {
		t.Errorf("Expected parent to be nil, but got '%v'", group.parent)
	}

	if len(group.middlewares) != 0 {
		t.Errorf("Expected middlewares to be empty, but got '%v'", group.middlewares)
	}

	if group.app != app {
		t.Errorf("Expected app to be '%v', but got '%v'", app, group.app)
	}
}

func TestGroup_AddRoute(t *testing.T) {
	app := NewApp()
	group := app.Group("/prefix")
	handlers := []HandlerFunc{func(c *Context) {}}
	group.AddRoute(MethodGet, "/path", handlers)

	searchHandlers, _ := app.router.findRoute(MethodGet, "/prefix/path")
	if reflect.ValueOf(searchHandlers[0]) != reflect.ValueOf(handlers[0]) {
		t.Errorf("Expected handlers to be '%v', but got '%v'", searchHandlers[0], handlers[0])
	}
}

func TestGroup_Use(t *testing.T) {
	app := NewApp()
	group := app.Group("/test")

	middleware := func(c *Context) {
		c.SetHeader("X-Test-Header", "123")
	}

	group.Use(middleware)

	group.Get("/header", func(c *Context) {
		header := c.Header("X-Test-Header")
		c.Text(StatusOK, header)
	})

	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod(MethodGet)
	ctx.Request.Header.SetRequestURI("/test/header")
	app.serveRequest(ctx)

	if string(ctx.Response.Header.Peek("X-Test-Header")) != "123" {
		t.Errorf("Expected header to be '123', but got '%s'", string(ctx.Response.Header.Peek("X-Test-Header")))
	}
}

func TestGroup_Get(t *testing.T) {
	app := NewApp()
	group := app.Group("/prefix")
	handlers := []HandlerFunc{func(c *Context) {}}
	group.Get("/path", handlers...)

	searchHandlers, _ := app.router.findRoute(MethodGet, "/prefix/path")
	if reflect.ValueOf(searchHandlers[0]) != reflect.ValueOf(handlers[0]) {
		t.Errorf("Expected handlers to be '%v', but got '%v'", searchHandlers[0], handlers[0])
	}
}

func TestGroup_Post(t *testing.T) {
	app := NewApp()
	group := app.Group("/prefix")
	handlers := []HandlerFunc{func(c *Context) {}}
	group.Post("/path", handlers...)

	searchHandlers, _ := app.router.findRoute(MethodPost, "/prefix/path")
	if reflect.ValueOf(searchHandlers[0]) != reflect.ValueOf(handlers[0]) {
		t.Errorf("Expected handlers to be '%v', but got '%v'", searchHandlers[0], handlers[0])
	}
}

func TestGroup_Put(t *testing.T) {
	app := NewApp()
	group := app.Group("/prefix")
	handlers := []HandlerFunc{func(c *Context) {}}
	group.Put("/path", handlers...)

	searchHandlers, _ := app.router.findRoute(MethodPut, "/prefix/path")
	if reflect.ValueOf(searchHandlers[0]) != reflect.ValueOf(handlers[0]) {
		t.Errorf("Expected handlers to be '%v', but got '%v'", searchHandlers[0], handlers[0])
	}
}

func TestGroup_Delete(t *testing.T) {
	app := NewApp()
	group := app.Group("/prefix")
	handlers := []HandlerFunc{func(c *Context) {}}
	group.Delete("/path", handlers...)

	searchHandlers, _ := app.router.findRoute(MethodDelete, "/prefix/path")
	if reflect.ValueOf(searchHandlers[0]) != reflect.ValueOf(handlers[0]) {
		t.Errorf("Expected handlers to be '%v', but got '%v'", searchHandlers[0], handlers[0])
	}
}

func TestGroup_Head(t *testing.T) {
	app := NewApp()
	group := app.Group("/prefix")
	handlers := []HandlerFunc{func(c *Context) {}}
	group.Head("/path", handlers...)

	searchHandlers, _ := app.router.findRoute(MethodHead, "/prefix/path")
	if reflect.ValueOf(searchHandlers[0]) != reflect.ValueOf(handlers[0]) {
		t.Errorf("Expected handlers to be '%v', but got '%v'", searchHandlers[0], handlers[0])
	}
}

func TestGroup_Patch(t *testing.T) {
	app := NewApp()
	group := app.Group("/prefix")
	handlers := []HandlerFunc{func(c *Context) {}}
	group.Patch("/path", handlers...)

	searchHandlers, _ := app.router.findRoute(MethodPatch, "/prefix/path")
	if reflect.ValueOf(searchHandlers[0]) != reflect.ValueOf(handlers[0]) {
		t.Errorf("Expected handlers to be '%v', but got '%v'", searchHandlers[0], handlers[0])
	}
}

func TestGroup_Options(t *testing.T) {
	app := NewApp()
	group := app.Group("/prefix")
	handlers := []HandlerFunc{func(c *Context) {}}
	group.Options("/path", handlers...)

	searchHandlers, _ := app.router.findRoute(MethodOptions, "/prefix/path")
	if reflect.ValueOf(searchHandlers[0]) != reflect.ValueOf(handlers[0]) {
		t.Errorf("Expected handlers to be '%v', but got '%v'", searchHandlers[0], handlers[0])
	}
}
