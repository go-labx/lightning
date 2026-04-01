package lightning

import (
	"encoding/json"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"text/template"

	"github.com/go-labx/lightlog"
	"github.com/valyala/fasthttp"
)

// HandlerFunc is a function type that represents the actual handler function for a route.
type HandlerFunc func(*Context)

// Middleware is an alias for HandlerFunc, representing middleware functions.
type Middleware = HandlerFunc

// Map is a shortcut for map[string]interface{}
type Map map[string]any

// Application is the main struct that holds the router, middlewares, and configuration.
type Application struct {
	Config        *Config
	router        *router
	middlewares   []HandlerFunc
	htmlTemplates *template.Template
	funcMap       template.FuncMap

	Logger *lightlog.ConsoleLogger

	server      *fasthttp.Server
	mu          sync.Mutex
	contextPool sync.Pool
}

// Config holds the configuration for the Application.
type Config struct {
	AppName            string
	JSONEncoder        JSONMarshal
	JSONDecoder        JSONUnmarshal
	NotFoundHandler    HandlerFunc
	EnableDebug        bool
	MaxRequestBodySize int64
}

// merge merges the given Config structs into the current Config.
func (c *Config) merge(configs ...*Config) *Config {
	for _, cfg := range configs {
		if cfg == nil {
			continue
		}
		if cfg.AppName != "" {
			c.AppName = cfg.AppName
		}
		if cfg.JSONEncoder != nil {
			c.JSONEncoder = cfg.JSONEncoder
		}
		if cfg.JSONDecoder != nil {
			c.JSONDecoder = cfg.JSONDecoder
		}
		if cfg.NotFoundHandler != nil {
			c.NotFoundHandler = cfg.NotFoundHandler
		}
		if cfg.EnableDebug {
			c.EnableDebug = cfg.EnableDebug
		}
		if cfg.MaxRequestBodySize > 0 {
			c.MaxRequestBodySize = cfg.MaxRequestBodySize
		}
	}
	return c
}

// defaultConfig returns a new Config with default values.
func defaultConfig() *Config {
	return &Config{
		AppName:         "lightning-app",
		JSONEncoder:     defaultJSONMarshal,
		JSONDecoder:     defaultJSONUnmarshal,
		NotFoundHandler: defaultNotFound,
		EnableDebug:     false,
	}
}

// defaultJSONMarshal is the default JSON marshaling function.
func defaultJSONMarshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

// defaultJSONUnmarshal is the default JSON unmarshaling function.
func defaultJSONUnmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

// NewApp returns a new instance of the Application struct.
func NewApp(c ...*Config) *Application {
	config := defaultConfig()
	config = config.merge(c...)

	app := &Application{
		Config: config,
		router: newRouter(),
		Logger: lightlog.NewConsoleLogger(config.AppName, lightlog.TRACE),
		contextPool: sync.Pool{
			New: func() interface{} {
				return &Context{index: -1}
			},
		},
	}
	app.middlewares = make([]HandlerFunc, 0)

	if app.Config.EnableDebug {
		app.Get("/__debug__/router_map", func(ctx *Context) {
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

// AddRoute adds a new route to the router.
// It composes the global middlewares, route-specific middlewares, and the actual handler function
// to form a single MiddlewareFunc, and then adds it to the router.
func (app *Application) AddRoute(method string, pattern string, handlers []HandlerFunc) {
	app.Logger.Debug(" %s\t-> %s", method, pattern)
	allHandlers := make([]HandlerFunc, 0)
	allHandlers = append(allHandlers, app.middlewares...)
	allHandlers = append(allHandlers, handlers...)

	app.router.addRoute(method, pattern, allHandlers)
}

// Get adds a new route with method "GET" to the router.
func (app *Application) Get(pattern string, handlers ...HandlerFunc) {
	app.AddRoute(MethodGet, pattern, handlers)
}

// Post adds a new route with method "POST" to the router.
func (app *Application) Post(pattern string, handlers ...HandlerFunc) {
	app.AddRoute(MethodPost, pattern, handlers)
}

// Put adds a new route with method "PUT" to the router.
func (app *Application) Put(pattern string, handlers ...HandlerFunc) {
	app.AddRoute(MethodPut, pattern, handlers)
}

// Delete adds a new route with method "DELETE" to the router.
func (app *Application) Delete(pattern string, handlers ...HandlerFunc) {
	app.AddRoute(MethodDelete, pattern, handlers)
}

// Head adds a new route with method "HEAD" to the router.
func (app *Application) Head(pattern string, handlers ...HandlerFunc) {
	app.AddRoute(MethodHead, pattern, handlers)
}

// Patch adds a new route with method "PATCH" to the router.
func (app *Application) Patch(pattern string, handlers ...HandlerFunc) {
	app.AddRoute(MethodPatch, pattern, handlers)
}

// Options adds a new route with method "OPTIONS" to the router.
func (app *Application) Options(pattern string, handlers ...HandlerFunc) {
	app.AddRoute(MethodOptions, pattern, handlers)
}

// Group returns a new instance of the Group struct with the given prefix.
func (app *Application) Group(prefix string) *Group {
	return newGroup(app, prefix)
}

// Static serves static files from the given root directory with the given prefix.
// If root is an absolute path, it is used directly. Otherwise, it is resolved relative
// to the executable's directory.
// If the file exists, it is served with a 200 status code.
// If the file does not exist, a 404 status code is returned with the text "Not Found".
func (app *Application) Static(root string, prefix string) {
	exPath := ""
	if filepath.IsAbs(root) {
		exPath = ""
	} else {
		ex, err := os.Executable()
		if err != nil {
			app.Logger.Warn("Failed to get executable path for static files: %v, using current directory", err)
			exPath = "."
		} else {
			exPath = filepath.Dir(ex)
		}
	}

	app.Get(path.Join(prefix, "/*"), func(ctx *Context) {
		fullFilePath := filepath.Join(exPath, root, strings.TrimPrefix(ctx.Path, prefix))

		if _, err := os.Stat(fullFilePath); !os.IsNotExist(err) {
			ctx.SkipFlush()
			ctx.SetStatus(StatusOK)
			ctx.ctx.SendFile(fullFilePath)
		} else {
			ctx.Text(StatusNotFound, "Not Found")
		}
	})
}

// SetFuncMap sets the funcMap in the Application struct to the funcMap passed in as an argument.
func (app *Application) SetFuncMap(funcMap template.FuncMap) {
	app.funcMap = funcMap
}

// LoadHTMLGlob loads HTML templates from a glob pattern and sets them in the Application struct.
// It uses the template.Must function to panic if there is an error parsing the templates.
// It also sets the funcMap in the Application struct to the funcMap passed in as an argument.
func (app *Application) LoadHTMLGlob(pattern string) {
	app.htmlTemplates = template.Must(template.New("").Funcs(app.funcMap).ParseGlob(pattern))
}

// RequestHandler returns a fasthttp.RequestHandler for the Application.
func (app *Application) RequestHandler() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		app.serveRequest(ctx)
	}
}

// serveRequest handles incoming HTTP requests by finding the matching route,
// creating a new Context, setting the route parameters, and executing the middleware chain.
func (app *Application) serveRequest(ctx *fasthttp.RequestCtx) {
	c := app.acquireContext(ctx)
	defer app.releaseContext(c)

	handlers, params := app.router.findRoute(c.Method, c.Path)

	if handlers == nil {
		handlers = append(app.middlewares, app.Config.NotFoundHandler)
	}
	c.setHandlers(handlers)
	c.setParams(params)
	c.setApp(app)

	c.Next()
	c.flush()
}

// acquireContext gets a Context from the pool and initializes it.
func (app *Application) acquireContext(ctx *fasthttp.RequestCtx) *Context {
	c := app.contextPool.Get().(*Context)
	c.ctx = ctx
	c.req = newRequest(ctx)
	c.res = newResponse(ctx)
	c.Method = c.req.method()
	c.Path = c.req.path()
	c.App = app
	c.data = contextData{}

	return c
}

// releaseContext resets and returns the Context to the pool.
func (app *Application) releaseContext(c *Context) {
	c.reset()
	app.contextPool.Put(c)
}

// Run starts the HTTP server and listens for incoming requests.
func (app *Application) Run(address ...string) error {
	addr := resolveAddress(address)
	app.Logger.Info("Starting application on address `%s` 🚀🚀🚀", addr)

	app.mu.Lock()
	app.server = &fasthttp.Server{
		Handler:            app.RequestHandler(),
		MaxRequestBodySize: int(app.Config.MaxRequestBodySize),
	}
	app.mu.Unlock()

	return app.server.ListenAndServe(addr)
}

// RunGraceful starts the HTTP server with graceful shutdown support.
// It listens for SIGINT and SIGTERM signals to trigger graceful shutdown.
// The shutdownTimeout specifies the maximum duration in seconds to wait for active connections to finish.
func (app *Application) RunGraceful(shutdownTimeout int, address ...string) error {
	addr := resolveAddress(address)
	app.Logger.Info("Starting application on address `%s` 🚀🚀🚀", addr)

	app.mu.Lock()
	app.server = &fasthttp.Server{
		Handler:            app.RequestHandler(),
		MaxRequestBodySize: int(app.Config.MaxRequestBodySize),
	}
	app.mu.Unlock()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- app.server.ListenAndServe(addr)
	}()

	select {
	case sig := <-quit:
		app.Logger.Info("Received signal %v, shutting down gracefully...", sig)
		if shutdownTimeout <= 0 {
			shutdownTimeout = 5
		}
		app.server.Shutdown()
		app.Logger.Info("Server stopped gracefully")
		return nil

	case err := <-serverErr:
		if err != nil {
			return err
		}
		return nil
	}
}

// Shutdown gracefully shuts down the server without interrupting active connections.
func (app *Application) Shutdown() {
	if app.server != nil {
		app.server.Shutdown()
	}
}
