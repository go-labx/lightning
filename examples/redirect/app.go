package main

import (
	"net/http"

	"github.com/go-labx/lightning"
)

func main() {
	app := lightning.NewApp()

	app.Get("/foo", func(ctx *lightning.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "/baz")
	})

	app.Get("/bar", func(ctx *lightning.Context) {
		ctx.Redirect(302, "/baz")
	})

	app.Get("/baz", func(ctx *lightning.Context) {
		ctx.JSON(http.StatusOK, map[string]string{
			"message": "pong",
		})
	})

	app.Run()
}
