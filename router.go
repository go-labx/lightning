package lightning

import (
	"fmt"
	"strings"
)

// trieNode represents a node in the trie data structure used by the router.
type trieNode struct {
	children map[string]*trieNode // A map of child nodes keyed by their string values
	isEnd    bool                 // boolean flag indicating whether the node marks the end of a route
	handlers []HandlerFunc        // `HandlerFunc` functions that handles requests for the node's route
	params   map[string]int       // a map of parameter names and their corresponding indices in the route pattern
	wildcard string               // a string representing the name of the wildcard parameter in the route pattern (if any)
}

// router represents the HTTP router.
type router struct {
	roots map[string]*trieNode
}

// newTrieNode creates a new instance of the `trieNode` struct with default values.
func newTrieNode() *trieNode {
	return &trieNode{
		children: make(map[string]*trieNode),
		isEnd:    false,
		handlers: make([]HandlerFunc, 0),
		params:   make(map[string]int),
		wildcard: "",
	}
}

// newRouter creates a new instance of the `router` struct with an empty `roots` map.
func newRouter() *router {
	return &router{
		roots: make(map[string]*trieNode),
	}
}

// addRoute adds a new route to the router.
func (r *router) addRoute(method string, pattern string, handlers []HandlerFunc) {
	if !isValidHTTPMethod(method) {
		panic(fmt.Sprintf("method `%s` is not a standard HTTP method", method))
	}
	root, ok := r.roots[method]
	if !ok {
		root = newTrieNode()
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
				root.children[":"] = newTrieNode()
			}
			root = root.children[":"]
		} else if part[0] == '*' {
			// wildcard
			name := part[1:]
			if root.children["*"] == nil {
				root.children["*"] = newTrieNode()
			}
			root = root.children["*"]
			root.wildcard = name
			break
		} else {
			// static
			if root.children[part] == nil {
				root.children[part] = newTrieNode()
			}
			root = root.children[part]
		}
	}

	root.isEnd = true        // mark the end of the route
	root.handlers = handlers // set the handlers for the route
	root.params = params     // set the parameters for the route
}

// findRoute is used to find the appropriate handler function for a given HTTP request method and URL pattern.
func (r *router) findRoute(method string, pattern string) ([]HandlerFunc, map[string]string) {
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

	return root.handlers, params
}
