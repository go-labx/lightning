package lightning

import "net/http"

// Group represents a group of routes with a common prefix and middleware.
type Group struct {
	app         *Application
	parent      *Group
	prefix      string
	middlewares []HandlerFunc
}

// newGroup creates a new Group with the given prefix and Application.
func newGroup(app *Application, prefix string) *Group {
	return &Group{
		app:         app,
		prefix:      prefix,
		parent:      nil,
		middlewares: []HandlerFunc{},
	}
}

// getFullPrefix returns the full prefix of the Group, including the prefixes of its ancestors.
func (g *Group) getFullPrefix() string {
	if g.parent == nil {
		return g.prefix
	}
	return g.parent.getFullPrefix() + g.prefix
}

// getMiddlewares returns the middleware functions of the Group and its ancestors.
func (g *Group) getMiddlewares() []HandlerFunc {
	if g.parent == nil {
		return g.middlewares
	}
	return append(g.parent.getMiddlewares(), g.middlewares...)
}

// Group creates a new Group with the given prefix and adds it as a child of the current Group.
func (g *Group) Group(prefix string) *Group {
	group := newGroup(g.app, prefix)
	group.parent = g
	return group
}

// AddRoute adds a new route to the Application with the given method, pattern, and handlers.
// The route's path is the full prefix of the Group concatenated with the given pattern.
// The route's handlers are the middleware functions of the Group and its ancestors concatenated with the given handlers.
func (g *Group) AddRoute(method string, pattern string, handlers []HandlerFunc) {
	handlers = append(g.getMiddlewares(), handlers...)
	path := g.getFullPrefix() + pattern
	g.app.AddRoute(method, path, handlers)
}

// Use adds the given middleware functions to the Group's middleware stack.
func (g *Group) Use(middlewares ...HandlerFunc) {
	g.middlewares = append(g.middlewares, middlewares...)
}

// Get adds a new GET route to the Application with the given pattern and handlers.
func (g *Group) Get(pattern string, handlers ...HandlerFunc) {
	g.AddRoute(http.MethodGet, pattern, handlers)
}

// Post adds a new POST route to the Application with the given pattern and handlers.
func (g *Group) Post(pattern string, handlers ...HandlerFunc) {
	g.AddRoute(http.MethodPost, pattern, handlers)
}

// Put adds a new PUT route to the Application with the given pattern and handlers.
func (g *Group) Put(pattern string, handlers ...HandlerFunc) {
	g.AddRoute(http.MethodPut, pattern, handlers)
}

// Delete adds a new DELETE route to the Application with the given pattern and handlers.
func (g *Group) Delete(pattern string, handlers ...HandlerFunc) {
	g.AddRoute(http.MethodDelete, pattern, handlers)
}

// Head adds a new HEAD route to the Application with the given pattern and handlers.
func (g *Group) Head(pattern string, handlers ...HandlerFunc) {
	g.AddRoute(http.MethodHead, pattern, handlers)
}

// Options adds a new OPTIONS route to the Application with the given pattern and handlers.
func (g *Group) Options(pattern string, handlers ...HandlerFunc) {
	g.AddRoute(http.MethodOptions, pattern, handlers)
}

// Patch adds a new PATCH route to the Application with the given pattern and handlers.
func (g *Group) Patch(pattern string, handlers ...HandlerFunc) {
	g.AddRoute(http.MethodPatch, pattern, handlers)
}
