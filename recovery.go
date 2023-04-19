package lightning

import (
	"fmt"
	"os"
	"runtime/debug"
)

// Recovery returns a middleware that recovers from panics and sends a 500 response with an error message.
func Recovery(handler ...HandlerFunc) Middleware {
	return func(ctx *Context) {
		defer func() {
			if r := recover(); r != nil {
				os.Stderr.WriteString(fmt.Sprintf("panic: %v\n%s\n", r, debug.Stack()))

				fn := defaultInternalServerError
				if len(handler) > 0 {
					fn = handler[0]
				}
				fn(ctx)
			}
		}()
		ctx.Next()
	}
}
