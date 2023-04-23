package main

import (
	"net/http"

	"github.com/go-labx/lightning"
)

func main() {
	app := lightning.NewApp()

	app.Get("/foo", func(ctx *lightning.Context) {
		// Redirect to /baz with a 301 status code
		ctx.Redirect(http.StatusMovedPermanently, "/baz")
	})

	app.Get("/bar", func(ctx *lightning.Context) {
		// Redirect to /baz with a 302 status code
		ctx.Redirect(302, "/baz")
	})

	app.Get("/baz", func(ctx *lightning.Context) {
		// Return a JSON response with a "message" key and "pong" value
		ctx.JSON(http.StatusOK, map[string]string{
			"message": "pong",
		})
	})

	app.Run()
}
