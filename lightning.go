package lightning

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"text/template"
	"time"

	"github.com/go-labx/lightlog"
)

// HandlerFunc is a function type that represents the actual handler function for a route.
type HandlerFunc func(*Context)
type Middleware = HandlerFunc

// Map is a shortcut for map[string]interface{}
type Map map[string]any

type Application struct {
	Config        *Config
	router        *router
	middlewares   []HandlerFunc
	htmlTemplates *template.Template
	funcMap       template.FuncMap

	Logger *lightlog.ConsoleLogger

	server      *http.Server
	mu          sync.Mutex
	contextPool sync.Pool
}

type Config struct {
	AppName            string
	JSONEncoder        JSONMarshal
	JSONDecoder        JSONUnmarshal
	NotFoundHandler    HandlerFunc
	EnableDebug        bool
	MaxRequestBodySize int64
}

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

func defaultConfig() *Config {
	return &Config{
		AppName:         "lightning-app",
		JSONEncoder:     json.Marshal,
		JSONDecoder:     json.Unmarshal,
		NotFoundHandler: defaultNotFound,
		EnableDebug:     false,
	}
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

// AddRoute is a function that adds a new route to the router.
// It composes the global middlewares, route-specific middlewares, and the actual handler function
// to form a single MiddlewareFunc, and then adds it to the router.
func (app *Application) AddRoute(method string, pattern string, handlers []HandlerFunc) {
	app.Logger.Debug(" %s\t-> %s", method, pattern)
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

// Static serves static files from the given root directory with the given prefix.
// If root is an absolute path, it is used directly. Otherwise, it is resolved relative
// to the executable's directory.
// If the file exists, it is served with a 200 status code using the http.ServeFile function.
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
			ctx.SetStatus(http.StatusOK)
			http.ServeFile(ctx.Res, ctx.Req, fullFilePath)
		} else {
			ctx.Text(http.StatusNotFound, http.StatusText(http.StatusNotFound))
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

// ServeHTTP is the function that handles HTTP requests.
// It finds the matching route, creates a new Context, sets the route parameters,
// and executes the MiddlewareFunc chain.
func (app *Application) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Apply max request body size limit
	if app.Config.MaxRequestBodySize > 0 {
		req.Body = http.MaxBytesReader(w, req.Body, app.Config.MaxRequestBodySize)
	}

	// Get context from pool
	ctx := app.acquireContext(w, req)
	defer app.releaseContext(ctx)

	// Find the matching route and set the handlers and paramsMap in the context
	handlers, params := app.router.findRoute(ctx.Method, ctx.Path)

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
	ctx.flush()
}

// acquireContext gets a Context from the pool and initializes it.
func (app *Application) acquireContext(w http.ResponseWriter, req *http.Request) *Context {
	r, err := newRequest(req)
	if err != nil {
		panic(err)
	}

	ctx := app.contextPool.Get().(*Context)
	ctx.Req = req
	ctx.Res = w
	ctx.req = r
	ctx.res = newResponse(req, w)
	ctx.Method = r.method
	ctx.Path = r.path
	ctx.App = app
	ctx.data = contextData{}

	return ctx
}

// releaseContext resets and returns the Context to the pool.
func (app *Application) releaseContext(ctx *Context) {
	ctx.reset()
	app.contextPool.Put(ctx)
}

// Run starts the HTTP server and listens for incoming requests.
func (app *Application) Run(address ...string) error {
	addr := resolveAddress(address)
	app.Logger.Info("Starting application on address `%s` 🚀🚀🚀", addr)

	app.mu.Lock()
	app.server = &http.Server{
		Addr:    addr,
		Handler: app,
	}
	app.mu.Unlock()

	return app.server.ListenAndServe()
}

// RunListener starts the HTTP server with an existing net.Listener.
func (app *Application) RunListener(listener net.Listener) error {
	app.mu.Lock()
	app.server = &http.Server{
		Handler: app,
	}
	app.mu.Unlock()

	return app.server.Serve(listener)
}

// Shutdown gracefully shuts down the server without interrupting active connections.
func (app *Application) Shutdown(ctx context.Context) error {
	app.mu.Lock()
	server := app.server
	app.mu.Unlock()

	if server == nil {
		return nil
	}
	return server.Shutdown(ctx)
}

// RunGraceful starts the HTTP server with graceful shutdown support.
// It listens for SIGINT and SIGTERM signals to trigger graceful shutdown.
// The shutdownTimeout specifies the maximum duration to wait for active connections to finish.
func (app *Application) RunGraceful(shutdownTimeout time.Duration, address ...string) error {
	addr := resolveAddress(address)
	app.Logger.Info("Starting application on address `%s` 🚀🚀🚀", addr)

	app.mu.Lock()
	app.server = &http.Server{
		Addr:    addr,
		Handler: app,
	}
	app.mu.Unlock()

	// Channel to listen for shutdown signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Channel to receive server errors
	serverErr := make(chan error, 1)

	go func() {
		serverErr <- app.server.ListenAndServe()
	}()

	// Wait for either interrupt signal or server error
	select {
	case sig := <-quit:
		app.Logger.Info("Received signal %v, shutting down gracefully...", sig)
		if shutdownTimeout <= 0 {
			shutdownTimeout = 5 * time.Second
		}
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		if err := app.Shutdown(ctx); err != nil {
			app.Logger.Error("Error during graceful shutdown: %v", err)
			return err
		}
		app.Logger.Info("Server stopped gracefully")
		return nil

	case err := <-serverErr:
		if err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	}
}
