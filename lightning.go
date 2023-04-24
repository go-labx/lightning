package lightning

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/go-labx/lightlog"
)

// HandlerFunc is a function type that represents the actual handler function for a route.
type HandlerFunc func(*Context)
type Middleware = HandlerFunc

// Map is a shortcut for map[string]interface{}
type Map map[string]any

type Application struct {
	Config      *Config
	router      *router
	middlewares []HandlerFunc

	Logger *lightlog.ConsoleLogger
}

var logger = lightlog.NewConsoleLogger("appLogger", lightlog.TRACE)

type Config struct {
	AppName         string
	JSONEncoder     JSONMarshal
	JSONDecoder     JSONUnmarshal
	NotFoundHandler HandlerFunc // Handler function for 404 Not Found error
	EnableDebug     bool
}

func (c *Config) merge(configs ...*Config) *Config {
	value := reflect.ValueOf(c).Elem()

	// iterate over all the configs passed in
	for _, config := range configs {
		v := reflect.ValueOf(config).Elem()
		t := reflect.TypeOf(config).Elem()

		// iterate over all the fields in the config
		for i := 0; i < t.NumField(); i++ {
			// if the field is not zero, set the value of the field in the current config to the value of the field in the passed in config
			if !v.Field(i).IsZero() {
				value.Field(i).Set(v.Field(i))
			}
		}
	}

	return c
}

func defaultConfig() *Config {
	return &Config{
		AppName:         "lightning-app",
		JSONEncoder:     json.Marshal,
		JSONDecoder:     json.Unmarshal,
		NotFoundHandler: defaultNotFound,
		EnableDebug:     true,
	}
}

// NewApp returns a new instance of the Application struct.
func NewApp(c ...*Config) *Application {
	config := defaultConfig()
	config = config.merge(c...)

	app := &Application{
		Config:      config,
		router:      newRouter(),
		middlewares: make([]HandlerFunc, 0),
		Logger:      logger,
	}

	if app.Config.EnableDebug {
		app.Get("/__debug__/router-map", func(ctx *Context) {
			ctx.JSON(200, app.router.Roots)
		})
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

// Use adds one or more Middlewares to the array of middlewares in the Application struct.
func (app *Application) Use(middlewares ...Middleware) {
	app.middlewares = append(app.middlewares, middlewares...)
}

// AddRoute is a function that adds a new route to the router.
// It composes the global middlewares, route-specific middlewares, and the actual handler function
// to form a single MiddlewareFunc, and then adds it to the router.
func (app *Application) AddRoute(method string, pattern string, handlers []HandlerFunc) {
	app.Logger.Debug("register route %s\t-> %s", method, pattern)
	allHandlers := make([]HandlerFunc, 0)
	allHandlers = append(allHandlers, app.middlewares...)
	allHandlers = append(allHandlers, handlers...)

	app.router.addRoute(method, pattern, allHandlers)
}

// The following functions are shortcuts for the addRoute function.
// They pre-fill the method parameter and call the addRoute function.

// Get adds a new route with method "GET" to the router.
func (app *Application) Get(pattern string, handlers ...HandlerFunc) {
	app.AddRoute("GET", pattern, handlers)
}

// Post adds a new route with method "POST" to the router.
func (app *Application) Post(pattern string, handlers ...HandlerFunc) {
	app.AddRoute("POST", pattern, handlers)
}

// Put adds a new route with method "PUT" to the router.
func (app *Application) Put(pattern string, handlers ...HandlerFunc) {
	app.AddRoute("PUT", pattern, handlers)
}

// Delete adds a new route with method "DELETE" to the router.
func (app *Application) Delete(pattern string, handlers ...HandlerFunc) {
	app.AddRoute("DELETE", pattern, handlers)
}

// Head adds a new route with method "HEAD" to the router.
func (app *Application) Head(pattern string, handlers ...HandlerFunc) {
	app.AddRoute("HEAD", pattern, handlers)
}

// Patch adds a new route with method "PATCH" to the router.
func (app *Application) Patch(pattern string, handlers ...HandlerFunc) {
	app.AddRoute("PATCH", pattern, handlers)
}

// Options adds a new route with method "OPTIONS" to the router.
func (app *Application) Options(pattern string, handlers ...HandlerFunc) {
	app.AddRoute("OPTIONS", pattern, handlers)
}

// Group returns a new instance of the Group struct with the given prefix.
func (app *Application) Group(prefix string) *Group {
	return newGroup(app, prefix)
}

// ServeHTTP is the function that handles HTTP requests.
// It finds the matching route, creates a new Context, sets the route parameters,
// and executes the MiddlewareFunc chain.
func (app *Application) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Create a new context
	ctx, err := NewContext(w, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer ctx.flush()

	// Find the matching route and set the handlers and paramsMap in the context
	handlers, params := app.router.findRoute(req.Method, req.URL.Path)
	// This check is necessary because if no matching route is found and the handlers slice is left empty,
	// the middleware chain will not be executed and the client will receive an empty response.
	// By appending the 404 handler function to the handlers slice,
	// we ensure that the middleware chain will always be executed, even if no matching route is found.
	if handlers == nil {
		handlers = append(app.middlewares, app.Config.NotFoundHandler)
	}
	ctx.setHandlers(handlers)
	ctx.setParams(params)
	ctx.setApp(app)

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
