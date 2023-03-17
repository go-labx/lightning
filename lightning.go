package lightning

import (
	"fmt"
	"log"
	"net/http"
)

// Compose takes multiple MiddlewareFuns and returns a single MiddlewareFunc.
// The returned MiddlewareFunc executes each MiddlewareFunc in the order they are passed,
// and then executes the `next` function.
func Compose(mw ...MiddlewareFunc) MiddlewareFunc {
	return func(c *Context, next Next) {
		var recursive func(i int)
		recursive = func(i int) {
			if i < len(mw) {
				// Execute the current middleware and pass the `next` function recursively.
				mw[i](c, func() {
					recursive(i + 1)
				})
			} else {
				// Execute the `next` function when all middlewares have been executed.
				next()
			}
		}
		recursive(0)
	}
}

// Next is a function type that represents the next middleware in the chain.
type Next func()

// HandlerFunc is a function type that represents the actual handler function for a route.
type HandlerFunc func(*Context)

// MiddlewareFunc is a function type that represents a middleware function.
type MiddlewareFunc func(ctx *Context, next Next)

type Application struct {
	router      *Router
	middlewares []MiddlewareFunc
}

// App returns a new instance of the Application struct.
func App() *Application {
	return &Application{
		router: NewRouter(),
	}
}

// Use adds one or more MiddlewareFuncs to the array of middlewares in the Application struct.
func (app *Application) Use(middlewares ...MiddlewareFunc) {
	app.middlewares = append(app.middlewares, middlewares...)
}

// AddRoute is a function that adds a new route to the Router.
// It composes the global middlewares, route-specific middlewares, and the actual handler function
// to form a single MiddlewareFunc, and then adds it to the Router.
func (app *Application) AddRoute(method string, pattern string, handler HandlerFunc, middlewares ...MiddlewareFunc) {
	log.Printf("register route %s %s ->", method, pattern)
	fn := Compose(Compose(app.middlewares...), Compose(middlewares...), func(ctx *Context, next Next) {
		handler(ctx)
	})

	app.router.AddRoute(method, pattern, func(ctx *Context) {
		fn(ctx, nil)
	})
}

// The following functions are shortcuts for the AddRoute function.
// They pre-fill the method parameter and call the AddRoute function.

func (app *Application) Get(pattern string, handler HandlerFunc, middlewares ...MiddlewareFunc) {
	app.AddRoute("GET", pattern, handler, middlewares...)
}

func (app *Application) Post(pattern string, handler HandlerFunc, middlewares ...MiddlewareFunc) {
	app.AddRoute("POST", pattern, handler, middlewares...)
}

func (app *Application) Put(pattern string, handler HandlerFunc, middlewares ...MiddlewareFunc) {
	app.AddRoute("PUT", pattern, handler, middlewares...)
}

func (app *Application) Delete(pattern string, handler HandlerFunc, middlewares ...MiddlewareFunc) {
	app.AddRoute("DELETE", pattern, handler, middlewares...)
}

func (app *Application) Head(pattern string, handler HandlerFunc, middlewares ...MiddlewareFunc) {
	app.AddRoute("HEAD", pattern, handler, middlewares...)
}

func (app *Application) Patch(pattern string, handler HandlerFunc, middlewares ...MiddlewareFunc) {
	app.AddRoute("PATCH", pattern, handler, middlewares...)
}

func (app *Application) Options(pattern string, handler HandlerFunc, middlewares ...MiddlewareFunc) {
	app.AddRoute("OPTIONS", pattern, handler, middlewares...)
}

// ServeHTTP is the function that handles HTTP requests.
// It finds the matching route, creates a new Context, sets the route parameters,
// and executes the MiddlewareFunc chain.
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

// Run starts the HTTP server and listens for incoming requests.
func (app *Application) Run(addr string) (err error) {
	return http.ListenAndServe(addr, app)
}
