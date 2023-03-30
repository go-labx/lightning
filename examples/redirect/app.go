package main

import (
	"github.com/go-labx/lightning"
	"net/http"
)

func main() {
	app := lightning.App()

	app.Get("/foo", func(ctx *lightning.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "/baz")
	})

	app.Get("/bar", func(ctx *lightning.Context) {
		ctx.Redirect(302, "/baz")
	})

	app.Get("/baz", func(ctx *lightning.Context) {
		ctx.JSON(map[string]string{
			"message": "pong",
		})
	})

	app.Run()
}
