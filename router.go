// Package lightning is a Go library that provides a lightweight, high-perform
package lightning

import (
	"net/http"
	"strings"
)

// TrieNode represents a node in the trie data structure used by the router.
type TrieNode struct {
	children map[string]*TrieNode // A map of child nodes keyed by their string values
	isEnd    bool                 // boolean flag indicating whether the node marks the end of a route
	handler  HandlerFunc          // a `HandlerFunc` function that handles requests for the node's route
	params   map[string]int       // a map of parameter names and their corresponding indices in the route pattern
	wildcard string               // a string representing the name of the wildcard parameter in the route pattern (if any)
}

// Router represents the HTTP router.
type Router struct {
	roots map[string]*TrieNode
}

// parsePattern splits a route pattern string into its individual parts.
func parsePattern(pattern string) []string {
	parts := strings.Split(pattern, "/")
	result := make([]string, 0)
	for _, part := range parts {
		if part != "" {
			result = append(result, part)
		}
	}

	return result
}

// NewTrieNode creates a new instance of the `TrieNode` struct with default values.
func NewTrieNode() *TrieNode {
	return &TrieNode{
		children: make(map[string]*TrieNode),
		isEnd:    false,
		handler:  nil,
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
func (r *Router) AddRoute(method string, pattern string, handler HandlerFunc) {
	if !isValidHTTPMethod(method) {
		return
	}
	root, ok := r.roots[method]
	if !ok {
		root = NewTrieNode()
		r.roots[method] = root
	}

	params := make(map[string]int)
	parts := parsePattern(pattern)
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

	root.isEnd = true
	root.handler = handler
	root.params = params
}

// FindRoute is used to find the appropriate handler function for a given HTTP request method and URL pattern.
func (r *Router) FindRoute(method string, pattern string) (HandlerFunc, map[string]string) {
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	params := make(map[string]string)
	values := make(map[int]string)

	parts := parsePattern(pattern)
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

	return root.handler, params
}

// HTTP Method Functions

// Get adds a GET route to the router.
func (r *Router) Get(path string, handler HandlerFunc) {
	r.AddRoute(http.MethodGet, path, handler)
}

// Post adds a POST route to the router.
func (r *Router) Post(path string, handler HandlerFunc) {
	r.AddRoute(http.MethodPost, path, handler)
}

// Put adds a PUT route to the router.
func (r *Router) Put(path string, handler HandlerFunc) {
	r.AddRoute(http.MethodPut, path, handler)
}

// Delete adds a DELETE route to the router.
func (r *Router) Delete(path string, handler HandlerFunc) {
	r.AddRoute(http.MethodDelete, path, handler)
}

// Head adds a HEAD route to the router.
func (r *Router) Head(path string, handler HandlerFunc) {
	r.AddRoute(http.MethodHead, path, handler)
}

// Patch adds a PATCH route to the router.
func (r *Router) Patch(path string, handler HandlerFunc) {
	r.AddRoute(http.MethodPatch, path, handler)
}

// Options adds an OPTIONS route to the router.
func (r *Router) Options(path string, handler HandlerFunc) {
	r.AddRoute(http.MethodOptions, path, handler)
}

func isValidHTTPMethod(method string) bool {
	switch method {
	case "GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS":
		return true
	}
	return false
}
