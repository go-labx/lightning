package lightning

import "net/http"

// Recovery returns a middleware that recovers from panics and sends a 500 response with an error message.
func Recovery() func(ctx *Context) {
	return func(ctx *Context) {
		defer func() {
			if err := recover(); err != nil {
				code := http.StatusInternalServerError

				switch err.(type) {
				case error:
					message := http.StatusText(code) + ": " + err.(error).Error()
					ctx.Text(code, message)
				default:
					ctx.Text(code, http.StatusText(code))
				}
			}
		}()
		ctx.Next()
	}
}
