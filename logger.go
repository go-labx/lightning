package lightning

import (
	"fmt"
	"time"
)

// Logger returns a middleware function that logs incoming requests
func Logger() func(ctx *Context) {
	return func(ctx *Context) {
		start := time.Now()

		ctx.Next()

		end := time.Now()
		elapsed := int(end.Sub(start) / time.Millisecond)
		datetime := start.Format("2006-01-02 15:04:05")

		fmt.Printf("%s %s %d %s %dms %s %s\n", datetime, ctx.RemoteAddr(), ctx.Status(), ctx.Method, elapsed, ctx.Path, ctx.UserAgent())
	}
}
