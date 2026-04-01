package lightning

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"text/template"
	"time"
)

func TestNewApp(t *testing.T) {
	app := NewApp()
	if app == nil {
		t.Error("NewApp returned nil")
	}
}

func TestDefaultApp(t *testing.T) {
	app := DefaultApp()

	// Assert that the Logger field is not nil
	if app.Logger == nil {
		t.Errorf("Expected Logger field to not be nil")
	}

	// Assert that the middlewares field has the expected length
	expectedMiddlewareLength := 2
	if len(app.middlewares) != expectedMiddlewareLength {
		t.Errorf("Expected middlewares field to have length %d, but got %d", expectedMiddlewareLength, len(app.middlewares))
	}
}

func TestUse(t *testing.T) {
	app := NewApp()

	// Define some middleware functions
	mw1 := func(c *Context) {}
	mw2 := func(c *Context) {}

	// Add the middleware functions to the app
	app.Use(mw1, mw2)

	// Check if the middleware functions were added correctly
	if len(app.middlewares) != 2 {
		t.Errorf("Expected 2 middleware functions, but got %d", len(app.middlewares))
	}
}

func TestAddRoute(t *testing.T) {
	app := NewApp()
	app.AddRoute("GET", "/test", []HandlerFunc{})
	route, _ := app.router.findRoute("GET", "/test")
	if route == nil {
		t.Errorf("Expected route to be added to router")
	}
}

func TestGetRoute(t *testing.T) {
	app := NewApp()
	app.Get("/test", func(c *Context) {
		c.Text(http.StatusOK, "Hello, World!")
	})

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.ServeHTTP)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "Hello, World!"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestPostRoute(t *testing.T) {
	app := NewApp()
	app.Post("/test", func(c *Context) {
		c.Text(http.StatusOK, "Hello, World!")
	})

	req, err := http.NewRequest("POST", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.ServeHTTP)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "Hello, World!"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestPutRoute(t *testing.T) {
	app := NewApp()
	app.Put("/test", func(c *Context) {
		c.Text(http.StatusOK, "Hello, World!")
	})

	req, err := http.NewRequest("PUT", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.ServeHTTP)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "Hello, World!"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestDeleteRoute(t *testing.T) {
	app := NewApp()
	app.Delete("/test", func(c *Context) {
		c.Text(http.StatusOK, "Hello, World!")
	})

	req, err := http.NewRequest("DELETE", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.ServeHTTP)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "Hello, World!"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestHeadRoute(t *testing.T) {
	app := NewApp()
	app.Head("/test", func(c *Context) {
		c.Text(http.StatusOK, "Hello, World!")
	})

	req, err := http.NewRequest("HEAD", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.ServeHTTP)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "Hello, World!"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestPatchRoute(t *testing.T) {
	app := NewApp()
	app.Patch("/test", func(c *Context) {
		c.Text(http.StatusOK, "Hello, World!")
	})

	req, err := http.NewRequest("PATCH", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.ServeHTTP)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "Hello, World!"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestOptionsRoute(t *testing.T) {
	app := NewApp()
	app.Options("/test", func(c *Context) {
		c.Text(http.StatusOK, "Hello, World!")
	})

	req, err := http.NewRequest("OPTIONS", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.ServeHTTP)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "Hello, World!"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestServeHTTP(t *testing.T) {
	app := NewApp()
	app.Get("/test", func(c *Context) {
		c.Text(http.StatusOK, "Hello, World!")
	})

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.ServeHTTP)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "Hello, World!"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestRun(t *testing.T) {
	app := NewApp()

	// Use dynamic port allocation to avoid port conflicts
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to create listener: %v", err)
	}
	addr := listener.Addr().String()

	go app.RunListener(listener)

	// Wait for the server to start
	time.Sleep(100 * time.Millisecond)

	// Send a GET request to the server
	resp, err := http.Get("http://" + addr)
	if err != nil {
		t.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Assert that the response status code is 404
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code %d, but got %d", http.StatusNotFound, resp.StatusCode)
	}
}

func TestConfigMerge(t *testing.T) {
	defaultCfg := &Config{
		AppName:         "default",
		JSONEncoder:     nil,
		NotFoundHandler: nil,
		EnableDebug:     false,
	}

	cfg := &Config{
		AppName:            "custom",
		EnableDebug:        true,
		MaxRequestBodySize: 1024,
	}

	merged := defaultCfg.merge(cfg)

	if merged.AppName != "custom" {
		t.Errorf("AppName = %v, want %v", merged.AppName, "custom")
	}
	if !merged.EnableDebug {
		t.Error("EnableDebug should be true")
	}
	if merged.MaxRequestBodySize != 1024 {
		t.Errorf("MaxRequestBodySize = %v, want %v", merged.MaxRequestBodySize, 1024)
	}
}

func TestConfigMergeWithNil(t *testing.T) {
	cfg := &Config{
		AppName: "test",
	}

	merged := cfg.merge(nil)
	if merged.AppName != "test" {
		t.Error("merge with nil should preserve original value")
	}

	merged = (&Config{}).merge(nil, nil, nil)
	if merged.AppName != "" {
		t.Error("merge with all nil should preserve zero values")
	}
}

func TestStatic(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "static_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	tmpFile := tmpDir + "/test.txt"
	if err := os.WriteFile(tmpFile, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	app := NewApp()
	app.Static(tmpDir, "/static")

	req := httptest.NewRequest("GET", "/static/test.txt", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
	if w.Body.String() != "hello" {
		t.Errorf("expected body 'hello', got %q", w.Body.String())
	}
}

func TestStaticWithAbsolutePath(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "static_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	tmpFile := tmpDir + "/test.txt"
	if err := os.WriteFile(tmpFile, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	absPath, _ := filepath.Abs(tmpDir)
	app := NewApp()
	app.Static(absPath, "/static")

	req := httptest.NewRequest("GET", "/static/test.txt", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestStaticNotFound(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "static_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	app := NewApp()
	app.Static(tmpDir, "/static")

	req := httptest.NewRequest("GET", "/static/nonexistent.txt", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestLoadHTMLGlob(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "templates")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	if err := os.WriteFile(tmpDir+"/test.html", []byte("<h1>{{.}}</h1>"), 0644); err != nil {
		t.Fatal(err)
	}

	app := NewApp()
	app.LoadHTMLGlob(tmpDir + "/*.html")

	if app.htmlTemplates == nil {
		t.Error("htmlTemplates should not be nil")
	}
}

func TestSetFuncMap(t *testing.T) {
	app := NewApp()
	funcMap := template.FuncMap{
		"upper": strings.ToUpper,
	}

	app.SetFuncMap(funcMap)

	if app.funcMap == nil {
		t.Error("funcMap should not be nil")
	}
	if app.funcMap["upper"] == nil {
		t.Error("funcMap should contain 'upper' function")
	}
}

func TestRunGracefulShutdown(t *testing.T) {
	app := NewApp()
	app.Get("/test", func(c *Context) {
		c.Text(http.StatusOK, "ok")
	})

	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal(err)
	}

	done := make(chan error, 1)
	go func() {
		done <- app.RunListener(listener)
	}()

	time.Sleep(50 * time.Millisecond)

	app.Shutdown(context.Background())

	select {
	case err := <-done:
		if err != nil && err != http.ErrServerClosed {
			t.Errorf("unexpected error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Error("timeout waiting for server shutdown")
	}
}

func TestShutdownWithoutServer(t *testing.T) {
	app := NewApp()
	err := app.Shutdown(context.Background())
	if err != nil {
		t.Errorf("Shutdown should not return error when server is nil: %v", err)
	}
}
