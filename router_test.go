package lightning

import (
	"net/http"
	"reflect"
	"testing"
)

func TestNewRouter(t *testing.T) {
	router := newRouter()

	if router == nil {
		t.Errorf("newRouter returned nil")
	}

	if len(router.Roots) != 0 {
		t.Errorf("newRouter did not initialize roots map correctly")
	}
}

func TestNewTrieNode(t *testing.T) {
	node := newTrieNode()

	if node == nil {
		t.Error("Failed to create a new trieNode instance")
	}

	if node.IsEnd != false {
		t.Errorf("Expected isEnd to be false, but got %v", node.IsEnd)
	}

	if len(node.handlers) != 0 {
		t.Errorf("newTrieNode did not initialize handlers correctly")
	}

	if node.Wildcard != "" {
		t.Errorf("Expected wildcard to be empty, but got %v", node.Wildcard)
	}

	if len(node.Params) != 0 {
		t.Errorf("Expected paramsMap to be empty, but got %v", node.Params)
	}
}

// Test Case 1: Adding a route with a static URL pattern and a valid handler function.
func TestAddRouteStaticPatternValidHandler(t *testing.T) {
	router := newRouter()
	method := "GET"
	pattern := "/home"
	handlers := []HandlerFunc{func(context *Context) {}}
	router.addRoute(method, pattern, handlers)

	// Assert that the route was added correctly.
	root := router.Roots[method]
	node, ok := root.Children["home"]
	if !ok {
		t.Errorf("expected route to be added, but wasn't")
	}
	if node == nil {
		t.Errorf("expected a non-nil node, but got nil")
	}
	if !node.IsEnd {
		t.Errorf("expected node to be an end node, but wasn't")
	}
	if len(node.handlers) == 0 {
		t.Errorf("expected node to have a non-nil handler, but got nil")
	}
	if node.Params == nil {
		t.Errorf("expected node to have non-nil paramsMap, but got nil")
	}
}

// Test Case 2: Adding a route with a parameterized URL pattern and a valid handler function.
func TestAddRouteParameterizedPatternValidHandler(t *testing.T) {
	router := newRouter()
	method := "GET"
	pattern := "/users/:id"
	handlers := []HandlerFunc{func(context *Context) {}}
	router.addRoute(method, pattern, handlers)

	// Assert that the route was added correctly.
	root := router.Roots[method]
	node, ok := root.Children["users"].Children[":"]
	if !ok {
		t.Errorf("Expected the route to be added, but it wasn't")
	}
	if node == nil {
		t.Errorf("Expected a non-nil node, but got nil")
	}
	if node.Params == nil {
		t.Errorf("Expected non-nil paramsMap, but got nil")
	}
	if node.Params["id"] != 1 {
		t.Errorf("Expected 1 parameter, but got %d", node.Params["id"])
	}
}

// Test Case 3: Adding a route with a wildcard URL pattern and a valid handler function.
func TestAddRouteWildcardPatternValidHandler(t *testing.T) {
	router := newRouter()
	method := "GET"
	pattern := "/users/*name"
	handlers := []HandlerFunc{func(context *Context) {}}
	router.addRoute(method, pattern, handlers)

	// Assert that the route was added correctly.
	root := router.Roots[method]
	node, ok := root.Children["users"]
	if !ok {
		t.Errorf("Expected the route to be added, but it wasn't")
	}
	node, ok = node.Children["*"]
	if node == nil {
		t.Errorf("Expected a non-nil node, but got nil")
	}
	if !node.IsEnd {
		t.Errorf("Expected node to be an end node, but wasn't")
	}
	if len(node.handlers) == 0 {
		t.Errorf("Expected node to have a non-nil handler, but got nil")
	}
	if node.Params == nil {
		t.Errorf("Expected node to have non-nil paramsMap, but got nil")
	}
	if node.Wildcard != "name" {
		t.Errorf("Expected node to have wildcard 'name', but got '%s'", node.Wildcard)
	}
}

// Test Case 4: Adding a route with an invalid method name.
func TestAddRouteInvalidMethod(t *testing.T) {
	router := newRouter()
	method := "INVALID"
	pattern := "/home"
	handlers := []HandlerFunc{func(context *Context) {}}

	// Use defer to capture the panic error
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected addRoute to panic with an error")
		}
	}()

	router.addRoute(method, pattern, handlers)
}

func TestRouter_FindRoute(t *testing.T) {
	// Create a new router instance
	router := newRouter()

	// Define the handler function for the test cases
	testHandler := func(ctx *Context) {}
	handlers := make([]HandlerFunc, 0)
	handlers = append(handlers, testHandler)

	// Add routes for test cases 3, 4, and 5
	router.addRoute(http.MethodGet, "/test", handlers)
	router.addRoute(http.MethodGet, "/users/:id", handlers)
	router.addRoute(http.MethodGet, "/files/*path", handlers)

	// Test case 1: invalid HTTP method
	if handlers, params := router.findRoute("INVALID_METHOD", "/test"); handlers != nil || params != nil {
		t.Errorf("Expected nil handler and paramsMap, but got handler %v and paramsMap %v", handlers, params)
	}

	// Test case 2: route does not exist
	if handlers, params := router.findRoute(http.MethodGet, "/invalid"); handlers != nil || params != nil {
		t.Errorf("Expected nil handler and paramsMap, but got handler %v and paramsMap %v", handlers, params)
	}

	// Test case 3: route exists with no parameters
	if handlers, params := router.findRoute(http.MethodGet, "/test"); reflect.ValueOf(handlers[0]).Pointer() != reflect.ValueOf(testHandler).Pointer() || len(params) != 0 {
		t.Errorf("Expected handler %v and empty paramsMap map, but got handler %v and paramsMap %v", "testHandler", handlers[0], params)
	}

	// Test case 4: route exists with parameters
	if handlers, params := router.findRoute(http.MethodGet, "/users/123"); reflect.ValueOf(handlers[0]).Pointer() != reflect.ValueOf(testHandler).Pointer() || len(params) != 1 || params["id"] != "123" {
		t.Errorf("Expected handler %v and paramsMap map {\"id\":\"123\"}, but got handler %v and paramsMap %v", "testHandler", handlers[0], params)
	}

	// Test case 5: route exists with wildcard parameter
	if handlers, params := router.findRoute(http.MethodGet, "/files/path/to/file.txt"); reflect.ValueOf(handlers[0]).Pointer() != reflect.ValueOf(testHandler).Pointer() || len(params) != 1 || params["path"] != "path/to/file.txt" {
		t.Errorf("Expected handler %v and paramsMap map {\"path\":\"path/to/file.txt\"}, but got handler %v and paramsMap %v", "testHandler", handlers[0], params)
	}
}
