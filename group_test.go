package lightning

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
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

	// Test that the group has the correct prefix
	if group.prefix != "/api" {
		t.Errorf("Expected prefix to be '/api', but got '%s'", group.prefix)
	}

	// Test that the group has the correct parent
	if group.parent != nil {
		t.Errorf("Expected parent to be nil, but got '%v'", group.parent)
	}

	// Test that the group has the correct middleware
	if len(group.middlewares) != 0 {
		t.Errorf("Expected middlewares to be empty, but got '%v'", group.middlewares)
	}

	// Test that the group has the correct application
	if group.app != app {
		t.Errorf("Expected app to be '%v', but got '%v'", app, group.app)
	}
}

func TestGroup_AddRoute(t *testing.T) {
	app := NewApp()
	group := app.Group("/prefix")
	handlers := []HandlerFunc{func(c *Context) {}}
	group.AddRoute(http.MethodGet, "/path", handlers)

	searchHandlers, _ := app.router.findRoute(http.MethodGet, "/prefix/path")
	if reflect.ValueOf(searchHandlers[0]) != reflect.ValueOf(handlers[0]) {
		t.Errorf("Expected handlers to be '%v', but got '%v'", searchHandlers[0], handlers[0])
	}
}

func TestGroup_Use(t *testing.T) {
	app := NewApp()
	group := app.Group("/test")

	// Define a middleware function that sets a custom header
	middleware := func(c *Context) {
		c.SetHeader("X-Test-Header", "123")
	}

	// Add the middleware function to the Group
	group.Use(middleware)

	// Define a route that returns the value of the custom header
	group.Get("/header", func(c *Context) {
		header := c.Header("X-Test-Header")
		c.Text(http.StatusOK, header)
	})

	// Send a request to the route using an HTTP client
	req, _ := http.NewRequest(http.MethodGet, "/test/header", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	// Verify that the response contains the custom header
	if w.Header().Get("X-Test-Header") != "123" {
		t.Errorf("Expected header to be '%v', but got '%v'", 123, w.Header().Get("X-Test-Header"))
	}
}

func TestGroup_Get(t *testing.T) {
	app := NewApp()
	group := app.Group("/prefix")
	handlers := []HandlerFunc{func(c *Context) {}}
	group.Get("/path", handlers...)

	searchHandlers, _ := app.router.findRoute(http.MethodGet, "/prefix/path")
	if reflect.ValueOf(searchHandlers[0]) != reflect.ValueOf(handlers[0]) {
		t.Errorf("Expected handlers to be '%v', but got '%v'", searchHandlers[0], handlers[0])
	}
}

func TestGroup_Post(t *testing.T) {
	app := NewApp()
	group := app.Group("/prefix")
	handlers := []HandlerFunc{func(c *Context) {}}
	group.Post("/path", handlers...)

	searchHandlers, _ := app.router.findRoute(http.MethodPost, "/prefix/path")
	if reflect.ValueOf(searchHandlers[0]) != reflect.ValueOf(handlers[0]) {
		t.Errorf("Expected handlers to be '%v', but got '%v'", searchHandlers[0], handlers[0])
	}
}

func TestGroup_Put(t *testing.T) {
	app := NewApp()
	group := app.Group("/prefix")
	handlers := []HandlerFunc{func(c *Context) {}}
	group.Put("/path", handlers...)

	searchHandlers, _ := app.router.findRoute(http.MethodPut, "/prefix/path")
	if reflect.ValueOf(searchHandlers[0]) != reflect.ValueOf(handlers[0]) {
		t.Errorf("Expected handlers to be '%v', but got '%v'", searchHandlers[0], handlers[0])
	}
}

func TestGroup_Delete(t *testing.T) {
	app := NewApp()
	group := app.Group("/prefix")
	handlers := []HandlerFunc{func(c *Context) {}}
	group.Delete("/path", handlers...)

	searchHandlers, _ := app.router.findRoute(http.MethodDelete, "/prefix/path")
	if reflect.ValueOf(searchHandlers[0]) != reflect.ValueOf(handlers[0]) {
		t.Errorf("Expected handlers to be '%v', but got '%v'", searchHandlers[0], handlers[0])
	}
}

func TestGroup_Head(t *testing.T) {
	app := NewApp()
	group := app.Group("/prefix")
	handlers := []HandlerFunc{func(c *Context) {}}
	group.Head("/path", handlers...)

	searchHandlers, _ := app.router.findRoute(http.MethodHead, "/prefix/path")
	if reflect.ValueOf(searchHandlers[0]) != reflect.ValueOf(handlers[0]) {
		t.Errorf("Expected handlers to be '%v', but got '%v'", searchHandlers[0], handlers[0])
	}
}

func TestGroup_Patch(t *testing.T) {
	app := NewApp()
	group := app.Group("/prefix")
	handlers := []HandlerFunc{func(c *Context) {}}
	group.Patch("/path", handlers...)

	searchHandlers, _ := app.router.findRoute(http.MethodPatch, "/prefix/path")
	if reflect.ValueOf(searchHandlers[0]) != reflect.ValueOf(handlers[0]) {
		t.Errorf("Expected handlers to be '%v', but got '%v'", searchHandlers[0], handlers[0])
	}
}

func TestGroup_Options(t *testing.T) {
	app := NewApp()
	group := app.Group("/prefix")
	handlers := []HandlerFunc{func(c *Context) {}}
	group.Options("/path", handlers...)

	searchHandlers, _ := app.router.findRoute(http.MethodOptions, "/prefix/path")
	if reflect.ValueOf(searchHandlers[0]) != reflect.ValueOf(handlers[0]) {
		t.Errorf("Expected handlers to be '%v', but got '%v'", searchHandlers[0], handlers[0])
	}
}
