package lightning

import (
	"fmt"
	"log"
	"net/http"
)

type HandlerFunc func(*Context)

type MiddlewareFunc func(HandlerFunc) HandlerFunc

type Application struct {
	router      *Router
	middlewares []MiddlewareFunc
}

func App() *Application {
	return &Application{
		router: NewRouter(),
	}
}

func (app *Application) addRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("register route %s %s ->", method, pattern)

	for _, m := range app.middlewares {
		handler = m(handler)
	}
	app.router.AddRoute(method, pattern, handler)
}

func (app *Application) Use(middlewares ...MiddlewareFunc) {
	app.middlewares = append(app.middlewares, middlewares...)
}

func (app *Application) GET(pattern string, handler HandlerFunc) {
	app.addRoute("GET", pattern, handler)
}

func (app *Application) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if handler, params := app.router.FindRoute(req.Method, req.URL.Path); handler != nil {
		ctx := NewContext(w, req)
		ctx.SetParams(params)

		handler(ctx)
	} else {
		_, err := fmt.Fprintf(w, "404 Not Found: %s\\n", req.URL)
		if err != nil {
			return
		}
	}
}

func (app *Application) Run(addr string) (err error) {
	return http.ListenAndServe(addr, app)
}
