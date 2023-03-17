package lightning

import (
	"net/http"
	"strings"
)

type TrieNode struct {
	children map[string]*TrieNode
	isEnd    bool
	handler  HandlerFunc
	params   map[string]int
	wildcard string
}

type Router struct {
	roots map[string]*TrieNode
}

func NewTrieNode() *TrieNode {
	return &TrieNode{
		children: make(map[string]*TrieNode),
		isEnd:    false,
		handler:  nil,
		params:   make(map[string]int),
		wildcard: "",
	}
}

func NewRouter() *Router {
	return &Router{
		roots: make(map[string]*TrieNode),
	}
}

func (r *Router) AddRoute(method string, pattern string, handler HandlerFunc) {
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

func (r *Router) FindHandler(method string, pattern string) (HandlerFunc, map[string]string) {
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

func (r *Router) Get(path string, handler HandlerFunc) {
	r.AddRoute(http.MethodGet, path, handler)
}

func (r *Router) Post(path string, handler HandlerFunc) {
	r.AddRoute(http.MethodPost, path, handler)
}

func (r *Router) Put(path string, handler HandlerFunc) {
	r.AddRoute(http.MethodPut, path, handler)
}

func (r *Router) Delete(path string, handler HandlerFunc) {
	r.AddRoute(http.MethodDelete, path, handler)
}

func (r *Router) Head(path string, handler HandlerFunc) {
	r.AddRoute(http.MethodHead, path, handler)
}

func (r *Router) Patch(path string, handler HandlerFunc) {
	r.AddRoute(http.MethodPatch, path, handler)
}

func (r *Router) Options(path string, handler HandlerFunc) {
	r.AddRoute(http.MethodOptions, path, handler)
}
