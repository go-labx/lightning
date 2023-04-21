package main

import (
	"github.com/go-labx/lightning"
)

func main() {
	app := lightning.DefaultApp()

	// Handle GET request to /success endpoint
	app.Get("/success", func(ctx *lightning.Context) {
		ctx.Success("hello world")
	})

	// Handle GET request to /fail endpoint
	app.Get("/fail", func(ctx *lightning.Context) {
		ctx.Fail(9999, "network error")
	})

	// Start the server
	app.Run()
}
