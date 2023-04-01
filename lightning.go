package lightning

import (
	"github.com/go-labx/lightlog"
	"net/http"
)

// HandlerFunc is a function type that represents the actual handler function for a route.
type HandlerFunc func(*Context)

// Map is a shortcut for map[string]interface{}
type Map map[string]any

type Application struct {
	router      *Router
	middlewares []HandlerFunc

	Logger                         *lightlog.ConsoleLogger
	NotFoundHandlerFunc            HandlerFunc // Handler function for 404 Not Found error
	InternalServerErrorHandlerFunc HandlerFunc // Handler function for 500 Internal Server Error
}

var logger = lightlog.NewConsoleLogger("appLogger", lightlog.TRACE)

// NewApp returns a new instance of the Application struct.
func NewApp() *Application {
	app := &Application{
		router:                         NewRouter(),
		middlewares:                    make([]HandlerFunc, 0),
		Logger:                         logger,
		NotFoundHandlerFunc:            defaultNotFound,
		InternalServerErrorHandlerFunc: defaultInternalServerError,
	}

	return app
}

// DefaultApp returns a new instance of the Application struct with default middlewares
func DefaultApp() *Application {
	app := NewApp()
	app.Use(Logger())
	app.Use(Recovery())

	return app
}

// Use adds one or more MiddlewareFuncs to the array of middlewares in the Application struct.
func (app *Application) Use(middlewares ...HandlerFunc) {
	app.middlewares = append(app.middlewares, middlewares...)
}

// AddRoute is a function that adds a new route to the Router.
// It composes the global middlewares, route-specific middlewares, and the actual handler function
// to form a single MiddlewareFunc, and then adds it to the Router.
func (app *Application) AddRoute(method string, pattern string, handlers []HandlerFunc) {
	app.Logger.Debug("register route %s\t-> %s", method, pattern)
	allHandlers := append(app.middlewares, handlers...)

	app.router.AddRoute(method, pattern, allHandlers)
}

// The following functions are shortcuts for the AddRoute function.
// They pre-fill the method parameter and call the AddRoute function.

// Get adds a new route with method "GET" to the Router.
func (app *Application) Get(pattern string, handlers ...HandlerFunc) {
	app.AddRoute("GET", pattern, handlers)
}

// Post adds a new route with method "POST" to the Router.
func (app *Application) Post(pattern string, handlers ...HandlerFunc) {
	app.AddRoute("POST", pattern, handlers)
}

// Put adds a new route with method "PUT" to the Router.
func (app *Application) Put(pattern string, handlers ...HandlerFunc) {
	app.AddRoute("PUT", pattern, handlers)
}

// Delete adds a new route with method "DELETE" to the Router.
func (app *Application) Delete(pattern string, handlers ...HandlerFunc) {
	app.AddRoute("DELETE", pattern, handlers)
}

// Head adds a new route with method "HEAD" to the Router.
func (app *Application) Head(pattern string, handlers ...HandlerFunc) {
	app.AddRoute("HEAD", pattern, handlers)
}

// Patch adds a new route with method "PATCH" to the Router.
func (app *Application) Patch(pattern string, handlers ...HandlerFunc) {
	app.AddRoute("PATCH", pattern, handlers)
}

// Options adds a new route with method "OPTIONS" to the Router.
func (app *Application) Options(pattern string, handlers ...HandlerFunc) {
	app.AddRoute("OPTIONS", pattern, handlers)
}

// Group returns a new instance of the Group struct with the given prefix.
func (app *Application) Group(prefix string) *Group {
	return NewGroup(app, prefix)
}

// ServeHTTP is the function that handles HTTP requests.
// It finds the matching route, creates a new Context, sets the route parameters,
// and executes the MiddlewareFunc chain.
func (app *Application) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Create a new context
	ctx, err := NewContext(w, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer ctx.Flush()

	// Find the matching route and set the handlers and params in the context
	handlers, params := app.router.FindRoute(req.Method, req.URL.Path)
	// This check is necessary because if no matching route is found and the handlers slice is left empty,
	// the middleware chain will not be executed and the client will receive an empty response.
	// By appending the 404 handler function to the handlers slice,
	// we ensure that the middleware chain will always be executed, even if no matching route is found.
	if handlers == nil {
		handlers = append(app.middlewares, app.NotFoundHandlerFunc)
	}
	ctx.SetHandlers(handlers)
	ctx.SetParams(params)

	// Execute the middleware chain
	ctx.Next()
}

// Run starts the HTTP server and listens for incoming requests.
func (app *Application) Run(address ...string) {
	addr := resolveAddress(address)
	app.Logger.Info("Starting application on address `%s` 🚀🚀🚀", addr)

	err := http.ListenAndServe(addr, app)
	if err != nil {
		panic(err)
	}
}
