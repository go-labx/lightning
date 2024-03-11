package lightning

import (
	"strings"
)

// node is a struct that represents a node in the trie
type node struct {
	// Pattern is the pattern of the node
	Pattern string `json:"pattern"`
	// Part is the part of the node
	Part string `json:"part"`
	// IsWild is a boolean that indicates whether the node is a wildcard
	IsWild bool `json:"isWild"`
	// Children is a slice of pointers to the children of the node
	Children []*node `json:"children,omitempty"`
	// handlers is a slice of HandlerFuncs that are associated with the node
	handlers []HandlerFunc
}

// matchChild returns the child node that matches the given part
func (n *node) matchChild(part string) *node {
	for _, child := range n.Children {
		if child.Part == part {
			return child
		}
	}
	return nil
}

// insert inserts a new node into the trie
func (n *node) insert(pattern string, parts []string, height int, handlers []HandlerFunc) {
	if len(parts) == height {
		n.Pattern = pattern
		n.handlers = handlers
		return
	}

	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{Part: part, IsWild: part[0] == ':' || part[0] == '*'}
		n.Children = append(n.Children, child)
	}
	child.insert(pattern, parts, height+1, handlers)
}

// search searches the trie for a node that matches the given parts
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height {
		if n.Pattern != "" {
			return n
		}
		return nil
	}

	part := parts[height]
	child := n.matchChild(part)

	// Attempt to match exact route first
	if child != nil {
		nextNode := child.search(parts, height+1)
		if nextNode != nil {
			return nextNode
		}
	}

	// Attempt to match wildcard route
	for _, child := range n.Children {
		if child.IsWild {
			nextNode := child.search(parts, height+1)
			if nextNode != nil {
				return nextNode
			}
		}
	}

	return nil
}

// router is a struct that represents a router
type router struct {
	// Roots is a map of HTTP methods to the root nodes of the trie
	Roots map[string]*node `json:"roots"`
}

// newRouter creates a new router
func newRouter() *router {
	return &router{
		Roots: make(map[string]*node, 0),
	}
}

// addRoute adds a new route to the router
func (r *router) addRoute(method string, pattern string, handlers []HandlerFunc) {
	parts := parsePattern(pattern)

	_, ok := r.Roots[method]
	if !ok {
		r.Roots[method] = &node{}
	}
	r.Roots[method].insert(pattern, parts, 0, handlers)
}

// findRoute finds the route that matches the given method and path
func (r *router) findRoute(method string, path string) ([]HandlerFunc, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.Roots[method]

	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0)

	if n != nil {
		parts := parsePattern(n.Pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n.handlers, params
	}

	return nil, nil
}
