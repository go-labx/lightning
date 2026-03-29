package lightning

import (
	"time"
)

// Logger returns a middleware function that logs incoming requests
func Logger() Middleware {
	return func(ctx *Context) {
		start := time.Now()

		ctx.Next()

		elapsed := time.Since(start)
		ctx.App.Logger.Info("%s %s %d %s %dms %s", ctx.RemoteAddr(), ctx.Method, ctx.Status(), ctx.Path, elapsed.Milliseconds(), ctx.UserAgent())
	}
}
