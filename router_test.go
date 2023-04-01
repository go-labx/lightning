package lightning

import (
	"net/http"
	"reflect"
	"testing"
)

func TestNewRouter(t *testing.T) {
	router := NewRouter()

	if router == nil {
		t.Errorf("NewRouter returned nil")
	}

	if len(router.roots) != 0 {
		t.Errorf("NewRouter did not initialize roots map correctly")
	}
}

func TestNewTrieNode(t *testing.T) {
	node := NewTrieNode()

	if node == nil {
		t.Error("Failed to create a new TrieNode instance")
	}

	if node.isEnd != false {
		t.Errorf("Expected isEnd to be false, but got %v", node.isEnd)
	}

	if len(node.handlers) != 0 {
		t.Errorf("NewTrieNode did not initialize handlers correctly")
	}

	if node.wildcard != "" {
		t.Errorf("Expected wildcard to be empty, but got %v", node.wildcard)
	}

	if len(node.params) != 0 {
		t.Errorf("Expected params to be empty, but got %v", node.params)
	}
}

// Test Case 1: Adding a route with a static URL pattern and a valid handler function.
func TestAddRouteStaticPatternValidHandler(t *testing.T) {
	router := NewRouter()
	method := "GET"
	pattern := "/home"
	handlers := []HandlerFunc{func(context *Context) {}}
	router.AddRoute(method, pattern, handlers)

	// Assert that the route was added correctly.
	root := router.roots[method]
	node, ok := root.children["home"]
	if !ok {
		t.Errorf("expected route to be added, but wasn't")
	}
	if node == nil {
		t.Errorf("expected a non-nil node, but got nil")
	}
	if !node.isEnd {
		t.Errorf("expected node to be an end node, but wasn't")
	}
	if len(node.handlers) == 0 {
		t.Errorf("expected node to have a non-nil handler, but got nil")
	}
	if node.params == nil {
		t.Errorf("expected node to have non-nil params, but got nil")
	}
}

// Test Case 2: Adding a route with a parameterized URL pattern and a valid handler function.
func TestAddRouteParameterizedPatternValidHandler(t *testing.T) {
	router := NewRouter()
	method := "GET"
	pattern := "/users/:id"
	handlers := []HandlerFunc{func(context *Context) {}}
	router.AddRoute(method, pattern, handlers)

	// Assert that the route was added correctly.
	root := router.roots[method]
	node, ok := root.children["users"].children[":"]
	if !ok {
		t.Errorf("Expected the route to be added, but it wasn't")
	}
	if node == nil {
		t.Errorf("Expected a non-nil node, but got nil")
	}
	if node.params == nil {
		t.Errorf("Expected non-nil params, but got nil")
	}
	if node.params["id"] != 1 {
		t.Errorf("Expected 1 parameter, but got %d", node.params["id"])
	}
}

// Test Case 3: Adding a route with a wildcard URL pattern and a valid handler function.
func TestAddRouteWildcardPatternValidHandler(t *testing.T) {
	router := NewRouter()
	method := "GET"
	pattern := "/users/*name"
	handlers := []HandlerFunc{func(context *Context) {}}
	router.AddRoute(method, pattern, handlers)

	// Assert that the route was added correctly.
	root := router.roots[method]
	node, ok := root.children["users"]
	if !ok {
		t.Errorf("Expected the route to be added, but it wasn't")
	}
	node, ok = node.children["*"]
	if node == nil {
		t.Errorf("Expected a non-nil node, but got nil")
	}
	if !node.isEnd {
		t.Errorf("Expected node to be an end node, but wasn't")
	}
	if len(node.handlers) == 0 {
		t.Errorf("Expected node to have a non-nil handler, but got nil")
	}
	if node.params == nil {
		t.Errorf("Expected node to have non-nil params, but got nil")
	}
	if node.wildcard != "name" {
		t.Errorf("Expected node to have wildcard 'name', but got '%s'", node.wildcard)
	}
}

// Test Case 4: Adding a route with an invalid method name.
func TestAddRouteInvalidMethod(t *testing.T) {
	router := NewRouter()
	method := "INVALID"
	pattern := "/home"
	handlers := []HandlerFunc{func(context *Context) {}}

	// Use defer to capture the panic error
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected AddRoute to panic with an error")
		}
	}()

	router.AddRoute(method, pattern, handlers)
}

func TestRouter_FindRoute(t *testing.T) {
	// Create a new router instance
	router := NewRouter()

	// Define the handler function for the test cases
	testHandler := func(ctx *Context) {}
	handlers := make([]HandlerFunc, 0)
	handlers = append(handlers, testHandler)

	// Add routes for test cases 3, 4, and 5
	router.AddRoute(http.MethodGet, "/test", handlers)
	router.AddRoute(http.MethodGet, "/users/:id", handlers)
	router.AddRoute(http.MethodGet, "/files/*path", handlers)

	// Test case 1: invalid HTTP method
	if handlers, params := router.FindRoute("INVALID_METHOD", "/test"); handlers != nil || params != nil {
		t.Errorf("Expected nil handler and params, but got handler %v and params %v", handlers, params)
	}

	// Test case 2: route does not exist
	if handlers, params := router.FindRoute(http.MethodGet, "/invalid"); handlers != nil || params != nil {
		t.Errorf("Expected nil handler and params, but got handler %v and params %v", handlers, params)
	}

	// Test case 3: route exists with no parameters
	if handlers, params := router.FindRoute(http.MethodGet, "/test"); reflect.ValueOf(handlers[0]).Pointer() != reflect.ValueOf(testHandler).Pointer() || len(params) != 0 {
		t.Errorf("Expected handler %v and empty params map, but got handler %v and params %v", "testHandler", handlers[0], params)
	}

	// Test case 4: route exists with parameters
	if handlers, params := router.FindRoute(http.MethodGet, "/users/123"); reflect.ValueOf(handlers[0]).Pointer() != reflect.ValueOf(testHandler).Pointer() || len(params) != 1 || params["id"] != "123" {
		t.Errorf("Expected handler %v and params map {\"id\":\"123\"}, but got handler %v and params %v", "testHandler", handlers[0], params)
	}

	// Test case 5: route exists with wildcard parameter
	if handlers, params := router.FindRoute(http.MethodGet, "/files/path/to/file.txt"); reflect.ValueOf(handlers[0]).Pointer() != reflect.ValueOf(testHandler).Pointer() || len(params) != 1 || params["path"] != "path/to/file.txt" {
		t.Errorf("Expected handler %v and params map {\"path\":\"path/to/file.txt\"}, but got handler %v and params %v", "testHandler", handlers[0], params)
	}
}
