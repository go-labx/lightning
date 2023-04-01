// Package lightning is a Go library that provides a lightweight, high-perform
package lightning

import (
	"fmt"
	"net/http"
	"strings"
)

// TrieNode represents a node in the trie data structure used by the router.
type TrieNode struct {
	children map[string]*TrieNode // A map of child nodes keyed by their string values
	isEnd    bool                 // boolean flag indicating whether the node marks the end of a route
	handlers []HandlerFunc        // `HandlerFunc` functions that handles requests for the node's route
	params   map[string]int       // a map of parameter names and their corresponding indices in the route pattern
	wildcard string               // a string representing the name of the wildcard parameter in the route pattern (if any)
}

// Router represents the HTTP router.
type Router struct {
	roots map[string]*TrieNode
}

// NewTrieNode creates a new instance of the `TrieNode` struct with default values.
func NewTrieNode() *TrieNode {
	return &TrieNode{
		children: make(map[string]*TrieNode),
		isEnd:    false,
		handlers: []HandlerFunc{},
		params:   make(map[string]int),
		wildcard: "",
	}
}

// NewRouter creates a new instance of the `Router` struct with an empty `roots` map.
func NewRouter() *Router {
	return &Router{
		roots: make(map[string]*TrieNode),
	}
}

// AddRoute adds a new route to the router.
func (r *Router) AddRoute(method string, pattern string, handlers []HandlerFunc) {
	if !isValidHTTPMethod(method) {
		panic(fmt.Sprintf("method `%s` is not a standard HTTP method", method))
	}
	root, ok := r.roots[method]
	if !ok {
		root = NewTrieNode()
		r.roots[method] = root
	}

	params := make(map[string]int)
	parts := ParsePattern(pattern)
	for i, part := range parts {
		if part[0] == ':' {
			// parameter
			name := part[1:]
			params[name] = i
			if root.children[":"] == nil {
				root.children[":"] = NewTrieNode()
			}
			root = root.children[":"]
		} else if part[0] == '*' {
			// wildcard
			name := part[1:]
			if root.children["*"] == nil {
				root.children["*"] = NewTrieNode()
			}
			root = root.children["*"]
			root.wildcard = name
			break
		} else {
			// static
			if root.children[part] == nil {
				root.children[part] = NewTrieNode()
			}
			root = root.children[part]
		}
	}

	root.isEnd = true // mark the end of the route
	root.handlers = handlers // set the handlers for the route
	root.params = params // set the parameters for the route
}

// FindRoute is used to find the appropriate handler function for a given HTTP request method and URL pattern.
func (r *Router) FindRoute(method string, pattern string) ([]HandlerFunc, map[string]string) {
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	params := make(map[string]string)
	values := make(map[int]string)

	parts := ParsePattern(pattern)
	for i, part := range parts {
		if root.children[part] != nil {
			root = root.children[part]
		} else if root.children[":"] != nil {
			root = root.children[":"]
			values[i] = part
		} else if root.children["*"] != nil {
			root = root.children["*"]
			if root.wildcard != "" {
				params[root.wildcard] = strings.Join(parts[i:], "/")
			}
			break
		} else {
			return nil, nil
		}
	}

	if !root.isEnd {
		return nil, nil
	}

	for name, index := range root.params {
		params[name] = values[index]
	}

	return root.handlers, params
}

// HTTP Method Functions

// Get adds a GET route to the router.
func (r *Router) Get(path string, handlers ...HandlerFunc) {
	r.AddRoute(http.MethodGet, path, handlers)
}

// Post adds a POST route to the router.
func (r *Router) Post(path string, handlers ...HandlerFunc) {
	r.AddRoute(http.MethodPost, path, handlers)
}

// Put adds a PUT route to the router.
func (r *Router) Put(path string, handlers ...HandlerFunc) {
	r.AddRoute(http.MethodPut, path, handlers)
}

// Delete adds a DELETE route to the router.
func (r *Router) Delete(path string, handlers ...HandlerFunc) {
	r.AddRoute(http.MethodDelete, path, handlers)
}

// Head adds a HEAD route to the router.
func (r *Router) Head(path string, handlers ...HandlerFunc) {
	r.AddRoute(http.MethodHead, path, handlers)
}

// Patch adds a PATCH route to the router.
func (r *Router) Patch(path string, handlers ...HandlerFunc) {
	r.AddRoute(http.MethodPatch, path, handlers)
}

// Options adds an OPTIONS route to the router.
func (r *Router) Options(path string, handlers ...HandlerFunc) {
	r.AddRoute(http.MethodOptions, path, handlers)
}

// isValidHTTPMethod checks if a given HTTP method is valid.
func isValidHTTPMethod(method string) bool {
	switch method {
	case "GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS":
		return true
	}
	return false
}
