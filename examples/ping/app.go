package main

import (
	"github.com/go-labx/lightning"
)

func main() {
	app := lightning.DefaultApp()

	app.Get("/ping", func(ctx *lightning.Context) {
		ctx.JSON(lightning.StatusOK, map[string]string{
			"message": "pong",
		})
	})

	app.Run()
}
