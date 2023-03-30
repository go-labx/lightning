package lightning

import "net/http"

type Group struct {
	app         *Application
	parent      *Group
	prefix      string
	middlewares []HandlerFunc
}

func NewGroup(app *Application, prefix string) *Group {
	return &Group{
		app:         app,
		prefix:      prefix,
		parent:      nil,
		middlewares: []HandlerFunc{},
	}
}

func (g *Group) getFulPrefix() string {
	if g.parent == nil {
		return g.prefix
	}
	return g.parent.getFulPrefix() + g.prefix
}

func (g *Group) getMiddlewares() []HandlerFunc {
	if g.parent == nil {
		return g.middlewares
	}
	return append(g.parent.getMiddlewares(), g.middlewares...)
}

func (g *Group) Group(prefix string) *Group {
	group := NewGroup(g.app, prefix)
	group.parent = g
	return group
}

func (g *Group) AddRoute(method string, pattern string, handlers []HandlerFunc) {
	handlers = append(g.getMiddlewares(), handlers...)
	path := g.getFulPrefix() + pattern
	g.app.AddRoute(method, path, handlers)
}

func (g *Group) Use(middlewares ...HandlerFunc) {
	g.middlewares = append(g.middlewares, middlewares...)
}

func (g *Group) Get(pattern string, handlers ...HandlerFunc) {
	g.AddRoute(http.MethodGet, pattern, handlers)
}

func (g *Group) Post(pattern string, handlers ...HandlerFunc) {
	g.AddRoute(http.MethodPost, pattern, handlers)
}

func (g *Group) Put(pattern string, handlers ...HandlerFunc) {
	g.AddRoute(http.MethodPut, pattern, handlers)
}

func (g *Group) Delete(pattern string, handlers ...HandlerFunc) {
	g.AddRoute(http.MethodDelete, pattern, handlers)
}

func (g *Group) Head(pattern string, handlers ...HandlerFunc) {
	g.AddRoute(http.MethodHead, pattern, handlers)
}

func (g *Group) Options(pattern string, handlers ...HandlerFunc) {
	g.AddRoute(http.MethodOptions, pattern, handlers)
}

func (g *Group) Patch(pattern string, handlers ...HandlerFunc) {
	g.AddRoute(http.MethodPatch, pattern, handlers)
}
