package lightning

import (
	"net/http"
	"net/http/httptest"
	"testing"
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

	// Assert that the NotFoundHandlerFunc field is not nil
	if app.NotFoundHandlerFunc == nil {
		t.Errorf("Expected NotFoundHandlerFunc field to not be nil")
	}

	// Assert that the InternalServerErrorHandlerFunc field is not nil
	if app.InternalServerErrorHandlerFunc == nil {
		t.Errorf("Expected InternalServerErrorHandlerFunc field to not be nil")
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
	go app.Run("localhost:9999")

	// Wait for the server to start
	time.Sleep(100 * time.Millisecond)

	// Send a GET request to the server
	resp, err := http.Get("http://localhost:9999")
	if err != nil {
		t.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Assert that the response status code is 404
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code %d, but got %d", http.StatusNotFound, resp.StatusCode)
	}
}
