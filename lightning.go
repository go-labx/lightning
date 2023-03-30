package lightning

import (
	"github.com/go-labx/lightlog"
	"net/http"
)

// HandlerFunc is a function type that represents the actual handler function for a route.
type HandlerFunc func(*Context)

type Application struct {
	router                         *Router
	middlewares                    []HandlerFunc
	logger                         *lightlog.ConsoleLogger
	NotFoundHandlerFunc            HandlerFunc
	InternalServerErrorHandlerFunc HandlerFunc
}

func DefaultNotFound(ctx *Context) {
	ctx.response.SetStatus(http.StatusNotFound)
	ctx.response.Text(http.StatusText(http.StatusNotFound))
}

func DefaultInternalServerError(ctx *Context) {
	ctx.response.SetStatus(http.StatusInternalServerError)
	ctx.response.Text(http.StatusText(http.StatusInternalServerError))
}

// NewApp returns a new instance of the Application struct.
func NewApp() *Application {
	app := &Application{
		router:                         NewRouter(),
		middlewares:                    make([]HandlerFunc, 0),
		logger:                         lightlog.NewConsoleLogger("logger", lightlog.TRACE),
		NotFoundHandlerFunc:            DefaultNotFound,
		InternalServerErrorHandlerFunc: DefaultInternalServerError,
	}

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
	app.logger.Trace("register route %s\t-> %s", method, pattern)
	allHandlers := append(app.middlewares, handlers...)

	app.router.AddRoute(method, pattern, allHandlers)
}

// The following functions are shortcuts for the AddRoute function.
// They pre-fill the method parameter and call the AddRoute function.

func (app *Application) Get(pattern string, handlers ...HandlerFunc) {
	app.AddRoute("GET", pattern, handlers)
}

func (app *Application) Post(pattern string, handlers ...HandlerFunc) {
	app.AddRoute("POST", pattern, handlers)
}

func (app *Application) Put(pattern string, handlers ...HandlerFunc) {
	app.AddRoute("PUT", pattern, handlers)
}

func (app *Application) Delete(pattern string, handlers ...HandlerFunc) {
	app.AddRoute("DELETE", pattern, handlers)
}

func (app *Application) Head(pattern string, handlers ...HandlerFunc) {
	app.AddRoute("HEAD", pattern, handlers)
}

func (app *Application) Patch(pattern string, handlers ...HandlerFunc) {
	app.AddRoute("PATCH", pattern, handlers)
}

func (app *Application) Options(pattern string, handlers ...HandlerFunc) {
	app.AddRoute("OPTIONS", pattern, handlers)
}

func (app *Application) Group(prefix string) *Group {
	return NewGroup(app, prefix)
}

// ServeHTTP is the function that handles HTTP requests.
// It finds the matching route, creates a new Context, sets the route parameters,
// and executes the MiddlewareFunc chain.
func (app *Application) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx, err := NewContext(w, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer ctx.Flush()

	handlers, params := app.router.FindRoute(req.Method, req.URL.Path)
	if handlers == nil {
		app.NotFoundHandlerFunc(ctx)
		return
	}

	ctx.SetHandlers(handlers)
	ctx.SetParams(params)
	ctx.Next()
}

// Run starts the HTTP server and listens for incoming requests.
func (app *Application) Run() {
	addr := "127.0.0.1:6789"
	app.logger.Info("Starting application on address `%s` 🚀🚀🚀", addr)

	err := http.ListenAndServe(addr, app)
	if err != nil {
		panic(err.Error())
	}
}
