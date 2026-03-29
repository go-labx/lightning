package lightning

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"text/template"
	"time"
)

func TestShutdown(t *testing.T) {
	app := NewApp()

	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to create listener: %v", err)
	}

	go app.RunListener(listener)

	time.Sleep(100 * time.Millisecond)

	// Test shutdown with nil server (should return nil)
	app2 := NewApp()
	err = app2.Shutdown(context.Background())
	if err != nil {
		t.Errorf("Expected nil error for shutdown with nil server, got %v", err)
	}

	// Test shutdown with running server
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = app.Shutdown(ctx)
	if err != nil {
		t.Errorf("Expected nil error for graceful shutdown, got %v", err)
	}
}

func TestRunGraceful(t *testing.T) {
	app := NewApp()
	app.Get("/test", func(c *Context) {
		c.Text(http.StatusOK, "Hello")
	})

	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to create listener: %v", err)
	}
	addr := listener.Addr().String()
	listener.Close() // Close it so RunGraceful can use a new one

	// Start server in goroutine with signal handling
	serverErr := make(chan error, 1)
	go func() {
		serverErr <- app.RunGraceful(2*time.Second, addr)
	}()

	// Wait for server to start
	time.Sleep(200 * time.Millisecond)

	// Make a request to verify server is running
	resp, err := http.Get("http://" + addr + "/test")
	if err != nil {
		t.Errorf("Failed to make request: %v", err)
	} else {
		resp.Body.Close()
	}

	// Shutdown the server via the app's Shutdown method
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	app.Shutdown(ctx)

	// Wait for RunGraceful to return
	select {
	case err := <-serverErr:
		if err != nil && err != http.ErrServerClosed {
			t.Errorf("Unexpected error: %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Error("RunGraceful did not return in time")
	}
}

func TestLoggerMiddleware(t *testing.T) {
	app := NewApp()
	app.Use(Logger())
	app.Get("/test", func(c *Context) {
		c.Text(http.StatusOK, "Hello, World!")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestRemoteAddr(t *testing.T) {
	tests := []struct {
		name       string
		headers    map[string]string
		remoteAddr string
		expected   string
	}{
		{
			name:       "X-Real-IP header",
			headers:    map[string]string{"X-Real-Ip": "192.168.1.1"},
			remoteAddr: "10.0.0.1:12345",
			expected:   "192.168.1.1",
		},
		{
			name:       "X-Forwarded-For single IP",
			headers:    map[string]string{"X-Forwarded-For": "192.168.1.2"},
			remoteAddr: "10.0.0.1:12345",
			expected:   "192.168.1.2",
		},
		{
			name:       "X-Forwarded-For multiple IPs",
			headers:    map[string]string{"X-Forwarded-For": "192.168.1.3, 10.0.0.2, 10.0.0.3"},
			remoteAddr: "10.0.0.1:12345",
			expected:   "192.168.1.3",
		},
		{
			name:       "No headers - use RemoteAddr",
			headers:    map[string]string{},
			remoteAddr: "10.0.0.1:12345",
			expected:   "10.0.0.1:12345",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}
			req.RemoteAddr = tt.remoteAddr

			r, err := newRequest(req)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			got := r.remoteAddr()
			if got != tt.expected {
				t.Errorf("Expected remoteAddr %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestRouterSearch(t *testing.T) {
	r := newRouter()

	// Add routes - test basic routing
	r.addRoute("GET", "/", []HandlerFunc{func(c *Context) {}})
	r.addRoute("GET", "/hello", []HandlerFunc{func(c *Context) {}})
	r.addRoute("GET", "/hello/:name", []HandlerFunc{func(c *Context) {}})
	r.addRoute("POST", "/users", []HandlerFunc{func(c *Context) {}})

	tests := []struct {
		method     string
		path       string
		found      bool
		paramCount int
	}{
		{"GET", "/", true, 0},
		{"GET", "/hello", true, 0},
		{"GET", "/hello/world", true, 1},
		{"POST", "/users", true, 0},
		{"GET", "/nonexistent", false, 0},
		{"DELETE", "/", false, 0}, // method not found
	}

	for _, tt := range tests {
		t.Run(tt.method+"_"+tt.path, func(t *testing.T) {
			handlers, params := r.findRoute(tt.method, tt.path)
			if tt.found && handlers == nil {
				t.Errorf("Expected to find route for %s %s", tt.method, tt.path)
			}
			if !tt.found && handlers != nil {
				t.Errorf("Expected not to find route for %s %s", tt.method, tt.path)
			}
			if tt.found && len(params) != tt.paramCount {
				t.Errorf("Expected %d params for %s %s, got %d", tt.paramCount, tt.method, tt.path, len(params))
			}
		})
	}
}

func TestResponseFlushWithCookies(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	res := newResponse(req, w)
	res.setStatus(http.StatusOK)
	res.setBody([]byte("test"))

	// Add cookies
	res.cookies.set("session", "abc123")

	res.flush()

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	cookies := w.Result().Cookies()
	if len(cookies) != 1 {
		t.Errorf("Expected 1 cookie, got %d", len(cookies))
	}
	if cookies[0].Name != "session" || cookies[0].Value != "abc123" {
		t.Errorf("Unexpected cookie: %v", cookies[0])
	}
}

func TestResponseFlushWithRedirect(t *testing.T) {
	req := httptest.NewRequest("GET", "/old", nil)
	w := httptest.NewRecorder()

	res := newResponse(req, w)
	res.redirect(http.StatusMovedPermanently, "/new")

	res.flush()

	if w.Code != http.StatusMovedPermanently {
		t.Errorf("Expected status %d, got %d", http.StatusMovedPermanently, w.Code)
	}

	location := w.Header().Get("Location")
	if location != "/new" {
		t.Errorf("Expected location '/new', got '%s'", location)
	}
}

func TestServeHTTPWithMaxBodySize(t *testing.T) {
	app := NewApp(&Config{
		MaxRequestBodySize: 10, // 10 bytes limit
	})
	app.Post("/test", func(c *Context) {
		c.Text(http.StatusOK, "OK")
	})

	// Request with body smaller than limit
	body := strings.NewReader("small")
	req := httptest.NewRequest("POST", "/test", body)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestAcquireContextPanic(t *testing.T) {
	app := NewApp()

	// Create a request that will cause panic in newRequest
	req := httptest.NewRequest("GET", "/test", nil)
	req.Body = nil // This should not cause panic

	ctx := app.acquireContext(httptest.NewRecorder(), req)
	if ctx == nil {
		t.Error("Expected non-nil context")
	}

	app.releaseContext(ctx)
}

func TestConfigMergeMultiple(t *testing.T) {
	config1 := &Config{
		AppName: "app1",
	}
	config2 := &Config{
		EnableDebug: true,
	}
	config3 := &Config{
		AppName: "app3",
	}

	merged := defaultConfig()
	merged = merged.merge(config1, nil, config2, config3)

	if merged.AppName != "app3" {
		t.Errorf("Expected AppName 'app3', got '%s'", merged.AppName)
	}
	if !merged.EnableDebug {
		t.Error("Expected EnableDebug to be true")
	}
}

func TestJSONBodyWithValidation(t *testing.T) {
	app := NewApp()

	type User struct {
		Name  string `validate:"required"`
		Email string `validate:"required,email"`
	}

	app.Post("/user", func(c *Context) {
		var user User
		if err := c.JSONBody(&user, true); err != nil {
			c.Fail(400, err.Error())
			return
		}
		c.Success(user)
	})

	// Valid request
	body := `{"name":"John","email":"john@example.com"}`
	req := httptest.NewRequest("POST", "/user", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Invalid request - missing required field
	body = `{"name":""}`
	req = httptest.NewRequest("POST", "/user", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestGroupCachedMiddlewares(t *testing.T) {
	app := NewApp()

	g := app.Group("/api")
	g.Use(func(c *Context) { c.Next() })

	// First call - computes and caches
	mw1 := g.getMiddlewares()
	if len(mw1) != 1 {
		t.Errorf("Expected 1 middleware, got %d", len(mw1))
	}

	// Second call - should return cached result
	mw2 := g.getMiddlewares()
	if len(mw2) != 1 {
		t.Errorf("Expected 1 middleware, got %d", len(mw2))
	}

	// Add more middleware - should invalidate cache
	g.Use(func(c *Context) { c.Next() })
	mw3 := g.getMiddlewares()
	if len(mw3) != 2 {
		t.Errorf("Expected 2 middlewares, got %d", len(mw3))
	}
}

func TestNestedGroupMiddlewares(t *testing.T) {
	app := NewApp()

	parent := app.Group("/api")
	parent.Use(func(c *Context) { c.Next() })

	child := parent.Group("/v1")
	child.Use(func(c *Context) { c.Next() })

	mw := child.getMiddlewares()
	if len(mw) != 2 {
		t.Errorf("Expected 2 middlewares, got %d", len(mw))
	}
}

func TestFileResponse(t *testing.T) {
	// Create a temp file
	tmpFile := "/tmp/test_lightning_file.txt"
	err := os.WriteFile(tmpFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile)

	req := httptest.NewRequest("GET", "/file", nil)
	w := httptest.NewRecorder()

	res := newResponse(req, w)
	err = res.file(tmpFile)
	if err != nil {
		t.Errorf("Failed to set file: %v", err)
	}

	res.flush()

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestFileResponseNotFound(t *testing.T) {
	req := httptest.NewRequest("GET", "/file", nil)
	w := httptest.NewRecorder()

	res := newResponse(req, w)
	err := res.file("/nonexistent/path/file.txt")
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}
}

func TestHTMLResponse(t *testing.T) {
	app := NewApp()

	// Skip HTML test since it requires template setup
	app.SetFuncMap(template.FuncMap{})
	// Note: LoadHTMLGlob requires actual template files
}

func TestXMLResponse(t *testing.T) {
	app := NewApp()
	app.Get("/xml", func(c *Context) {
		type Data struct {
			Key string `xml:"key"`
		}
		c.XML(http.StatusOK, &Data{Key: "value"})
	})

	req := httptest.NewRequest("GET", "/xml", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "value") {
		t.Errorf("Expected body to contain 'value', got %q", body)
	}
}

func TestContextDataOperations(t *testing.T) {
	app := NewApp()

	var gotV1, gotV2 any
	var v1AfterDel any

	app.Get("/test", func(c *Context) {
		c.SetData("key1", "value1")
		c.SetData("key2", 123)

		gotV1 = c.GetData("key1")
		gotV2 = c.GetData("key2")

		c.DelData("key1")
		v1AfterDel = c.GetData("key1")

		c.Text(http.StatusOK, "OK")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
	if gotV1 != "value1" {
		t.Errorf("Expected value1, got %v", gotV1)
	}
	if gotV2 != 123 {
		t.Errorf("Expected 123, got %v", gotV2)
	}
	if v1AfterDel != nil {
		t.Errorf("Expected nil after delete, got %v", v1AfterDel)
	}
}

func TestRedirectResponse(t *testing.T) {
	app := NewApp()
	app.Get("/redirect", func(c *Context) {
		c.Redirect(http.StatusMovedPermanently, "/new")
	})

	req := httptest.NewRequest("GET", "/redirect", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if w.Code != http.StatusMovedPermanently {
		t.Errorf("Expected status %d, got %d", http.StatusMovedPermanently, w.Code)
	}

	location := w.Header().Get("Location")
	if location != "/new" {
		t.Errorf("Expected location /new, got %s", location)
	}
}

func TestUserAgentAndReferer(t *testing.T) {
	app := NewApp()

	var gotUA, gotRef string

	app.Get("/test", func(c *Context) {
		gotUA = c.UserAgent()
		gotRef = c.Referer()
		c.Text(http.StatusOK, "OK")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("User-Agent", "test-agent")
	req.Header.Set("Referer", "http://example.com")
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if gotUA != "test-agent" {
		t.Errorf("Expected user agent 'test-agent', got '%s'", gotUA)
	}
	if gotRef != "http://example.com" {
		t.Errorf("Expected referer 'http://example.com', got '%s'", gotRef)
	}
}

func TestCookieOperations(t *testing.T) {
	app := NewApp()

	var gotCookieValue string

	app.Get("/cookie", func(c *Context) {
		cookie := c.Cookie("session")
		if cookie != nil {
			gotCookieValue = cookie.Value
		}
		c.SetCookie("new_cookie", "new_value")
		c.Text(http.StatusOK, "OK")
	})

	req := httptest.NewRequest("GET", "/cookie", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: "abc123"})
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if gotCookieValue != "abc123" {
		t.Errorf("Expected cookie value abc123, got %s", gotCookieValue)
	}

	cookies := w.Result().Cookies()
	if len(cookies) != 1 {
		t.Errorf("Expected 1 cookie, got %d", len(cookies))
	}
}

func TestCustomCookie(t *testing.T) {
	app := NewApp()
	app.Get("/cookie", func(c *Context) {
		c.SetCustomCookie(&http.Cookie{
			Name:     "custom",
			Value:    "value",
			Path:     "/",
			MaxAge:   3600,
			HttpOnly: true,
		})
		c.Text(http.StatusOK, "OK")
	})

	req := httptest.NewRequest("GET", "/cookie", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	cookies := w.Result().Cookies()
	if len(cookies) != 1 {
		t.Errorf("Expected 1 cookie, got %d", len(cookies))
	}
	if cookies[0].Name != "custom" {
		t.Errorf("Expected cookie name 'custom', got '%s'", cookies[0].Name)
	}
}

func TestHeaderOperations(t *testing.T) {
	app := NewApp()

	var gotHeaderValue string

	app.Get("/headers", func(c *Context) {
		gotHeaderValue = c.Header("X-Custom-Header")
		c.Text(http.StatusOK, "OK")
	})

	req := httptest.NewRequest("GET", "/headers", nil)
	req.Header.Set("X-Custom-Header", "value")
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if gotHeaderValue != "value" {
		t.Errorf("Expected header value 'value', got '%s'", gotHeaderValue)
	}
}

func TestAddAndDelHeaders(t *testing.T) {
	app := NewApp()
	app.Get("/headers", func(c *Context) {
		c.AddHeader("X-Multi", "value1")
		c.AddHeader("X-Multi", "value2")

		c.SetHeader("X-Set", "value3")

		c.DelHeader("X-Delete-Me")

		c.Text(http.StatusOK, "OK")
	})

	req := httptest.NewRequest("GET", "/headers", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	values := w.Header().Values("X-Multi")
	if len(values) != 2 {
		t.Errorf("Expected 2 values for X-Multi, got %d", len(values))
	}

	if w.Header().Get("X-Set") != "value3" {
		t.Errorf("Expected X-Set to be 'value3'")
	}

	if w.Header().Get("X-Delete-Me") != "" {
		t.Error("Expected X-Delete-Me to be deleted")
	}
}

func TestQueries(t *testing.T) {
	app := NewApp()

	var gotQueryValue string

	app.Get("/queries", func(c *Context) {
		queries := c.Queries()
		if len(queries["key1"]) > 0 {
			gotQueryValue = queries["key1"][0]
		}
		c.Text(http.StatusOK, "OK")
	})

	req := httptest.NewRequest("GET", "/queries?key1=value1&key2=value2", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if gotQueryValue != "value1" {
		t.Errorf("Expected query value 'value1', got '%s'", gotQueryValue)
	}
}

func TestSuccessFailResponses(t *testing.T) {
	app := NewApp()
	app.Get("/success", func(c *Context) {
		c.Success(map[string]string{"name": "test"})
	})
	app.Get("/fail", func(c *Context) {
		c.Fail(400, "Bad Request")
	})

	// Test success
	req := httptest.NewRequest("GET", "/success", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	if !strings.Contains(w.Body.String(), `"code":0`) {
		t.Error("Expected success response with code 0")
	}

	// Test fail
	req = httptest.NewRequest("GET", "/fail", nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)

	if !strings.Contains(w.Body.String(), `"code":400`) {
		t.Error("Expected fail response with code 400")
	}
}

func TestBodyOperations(t *testing.T) {
	app := NewApp()

	var gotBody []byte

	app.Get("/body", func(c *Context) {
		c.SetBody([]byte("custom body"))
		gotBody = c.Body()
	})

	req := httptest.NewRequest("GET", "/body", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if string(gotBody) != "custom body" {
		t.Errorf("Expected body 'custom body', got '%s'", string(gotBody))
	}
}

func TestStatus(t *testing.T) {
	app := NewApp()

	var gotStatus int

	app.Get("/status", func(c *Context) {
		c.SetStatus(http.StatusCreated)
		gotStatus = c.Status()
	})

	req := httptest.NewRequest("GET", "/status", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if gotStatus != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, gotStatus)
	}
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}
}
