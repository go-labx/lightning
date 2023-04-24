package lightning

import (
	"fmt"
	"strings"
)

// trieNode represents a node in the trie data structure used by the router.
type trieNode struct {
	Children map[string]*trieNode `json:"children"` // A map of child nodes keyed by their string values
	IsEnd    bool                 `json:"isEnd"`    // boolean flag indicating whether the node marks the end of a route
	handlers []HandlerFunc        // `HandlerFunc` functions that handles requests for the node's route
	Params   map[string]int       `json:"params"`   // a map of parameter names and their corresponding indices in the route pattern
	Wildcard string               `json:"wildcard"` // a string representing the name of the wildcard parameter in the route pattern (if any)
}

// router represents the HTTP router.
type router struct {
	Roots map[string]*trieNode `json:"roots"`
}

// newTrieNode creates a new instance of the `trieNode` struct with default values.
func newTrieNode() *trieNode {
	return &trieNode{
		Children: make(map[string]*trieNode),
		IsEnd:    false,
		handlers: make([]HandlerFunc, 0),
		Params:   make(map[string]int),
		Wildcard: "",
	}
}

// newRouter creates a new instance of the `router` struct with an empty `roots` map.
func newRouter() *router {
	return &router{
		Roots: make(map[string]*trieNode),
	}
}

// addRoute adds a new route to the router.
func (r *router) addRoute(method string, pattern string, handlers []HandlerFunc) {
	if !isValidHTTPMethod(method) {
		panic(fmt.Sprintf("method `%s` is not a standard HTTP method", method))
	}
	root, ok := r.Roots[method]
	if !ok {
		root = newTrieNode()
		r.Roots[method] = root
	}

	params := make(map[string]int)
	parts := parsePattern(pattern)
	for i, part := range parts {
		if part[0] == ':' {
			// parameter
			name := part[1:]
			params[name] = i
			if root.Children[":"] == nil {
				root.Children[":"] = newTrieNode()
			}
			root = root.Children[":"]
		} else if part[0] == '*' {
			// wildcard
			name := part[1:]
			if root.Children["*"] == nil {
				root.Children["*"] = newTrieNode()
			}
			root = root.Children["*"]
			root.Wildcard = name
			break
		} else {
			// static
			if root.Children[part] == nil {
				root.Children[part] = newTrieNode()
			}
			root = root.Children[part]
		}
	}

	root.IsEnd = true        // mark the end of the route
	root.handlers = handlers // set the handlers for the route
	root.Params = params     // set the parameters for the route
}

// findRoute is used to find the appropriate handler function for a given HTTP request method and URL pattern.
func (r *router) findRoute(method string, pattern string) ([]HandlerFunc, map[string]string) {
	root, ok := r.Roots[method]
	if !ok {
		return nil, nil
	}
	params := make(map[string]string)
	values := make(map[int]string)

	parts := parsePattern(pattern)
	for i, part := range parts {
		if root.Children[part] != nil {
			root = root.Children[part]
		} else if root.Children[":"] != nil {
			root = root.Children[":"]
			values[i] = part
		} else if root.Children["*"] != nil {
			root = root.Children["*"]
			if root.Wildcard != "" {
				params[root.Wildcard] = strings.Join(parts[i:], "/")
			}
			break
		} else {
			return nil, nil
		}
	}

	if !root.IsEnd {
		return nil, nil
	}

	for name, index := range root.Params {
		params[name] = values[index]
	}

	return root.handlers, params
}
