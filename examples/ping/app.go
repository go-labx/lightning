package main

import (
	"github.com/go-labx/lightning"
	"net/http"
)

func main() {
	app := lightning.DefaultApp()

	app.Get("/ping", func(ctx *lightning.Context) {
		ctx.JSON(http.StatusOK, lightning.Map{
			"message": "pong",
		})
	})

	app.Run()
}
