package main

import (
	"github.com/go-labx/lightning"
)

func main() {
	app := lightning.NewApp()

	app.Get("/foo", func(ctx *lightning.Context) {
		ctx.Redirect(lightning.StatusMovedPermanently, "/baz")
	})

	app.Get("/bar", func(ctx *lightning.Context) {
		ctx.Redirect(lightning.StatusFound, "/baz")
	})

	app.Get("/baz", func(ctx *lightning.Context) {
		ctx.JSON(lightning.StatusOK, map[string]string{
			"message": "pong",
		})
	})

	app.Run()
}
