package lightning

import "net/http"

func Recovery() func(ctx *Context) {
	return func(ctx *Context) {
		defer func() {
			if err := recover(); err != nil {
				switch err.(type) {
				case error:
					message := http.StatusText(http.StatusInternalServerError) + ": " + err.(error).Error()
					ctx.Text(http.StatusInternalServerError, message)
				default:
					ctx.Text(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
				}
			}
		}()
		ctx.Next()
	}
}
