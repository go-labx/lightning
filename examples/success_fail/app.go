package main

import (
	"github.com/go-labx/lightning"
)

func main() {
	app := lightning.DefaultApp()

	app.Get("/success", func(ctx *lightning.Context) {
		ctx.Success("hello world")
	})

	app.Get("/fail", func(ctx *lightning.Context) {
		ctx.Fail(9999, "network error")
	})

	app.Run()
}
