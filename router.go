package lightning

import (
	"net/http"
	"strings"
)

type TrieNode struct {
	children map[string]*TrieNode
	isLeaf   bool
	handler  HandlerFunc
}

type Router struct {
	roots map[string]*TrieNode
}

func NewTrieNode() *TrieNode {
	return &TrieNode{
		children: make(map[string]*TrieNode),
		isLeaf:   false,
		handler:  nil,
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

	parts := parsePattern(pattern)
	for _, part := range parts {
		if root.children[part] == nil {
			root.children[part] = NewTrieNode()
		}
		root = root.children[part]
	}

	root.isLeaf = true
	root.handler = handler
}

func (r *Router) FindHandler(method string, pattern string) HandlerFunc {
	root, ok := r.roots[method]
	if !ok {
		return nil
	}

	parts := parsePattern(pattern)
	for _, part := range parts {
		if root.children[part] != nil {
			root = root.children[part]
		}
	}
	return root.handler
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
