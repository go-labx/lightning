package lightning

import (
	"github.com/go-labx/lightlog"
	"time"
)

// Logger returns a middleware function that logs incoming requests
func Logger() func(ctx *Context) {
	logger := lightlog.NewConsoleLogger("accessLogger", lightlog.TRACE)

	return func(ctx *Context) {
		start := time.Now()
		method := ctx.Method
		path := ctx.Path

		ctx.Next()

		elapsed := int(time.Now().Sub(start) / time.Microsecond)
		logger.Trace("%s %d ---> %s %dms", method, ctx.Status(), path, elapsed)
	}
}
