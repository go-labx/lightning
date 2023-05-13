package lightning

import "strings"

type node struct {
	Pattern  string  `json:"pattern"`
	Part     string  `json:"part"`
	IsWild   bool    `json:"isWild"`
	Children []*node `json:"children,omitempty"`
	handlers []HandlerFunc
}

func (n *node) matchChild(part string) *node {
	for _, child := range n.Children {
		if child.Part == part || child.IsWild {
			return child
		}
	}
	return nil
}

func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.Children {
		if child.Part == part || child.IsWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

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

func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.Part, "*") {
		if n.Pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}

type router struct {
	Roots map[string]*node `json:"roots"`
}

func newRouter() *router {
	return &router{
		Roots: make(map[string]*node, 0),
	}
}

func (r *router) addRoute(method string, pattern string, handlers []HandlerFunc) {
	parts := parsePattern(pattern)

	_, ok := r.Roots[method]
	if !ok {
		r.Roots[method] = &node{}
	}
	r.Roots[method].insert(pattern, parts, 0, handlers)
}

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
