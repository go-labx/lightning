package lightning

// Helmet returns a middleware that sets common security response headers.
func Helmet() Middleware {
	return func(ctx *Context) {
		ctx.SetHeader("X-Content-Type-Options", "nosniff")
		ctx.SetHeader("X-Frame-Options", "DENY")
		ctx.SetHeader("X-XSS-Protection", "1; mode=block")
		ctx.SetHeader("Referrer-Policy", "strict-origin-when-cross-origin")
		ctx.Next()
	}
}
