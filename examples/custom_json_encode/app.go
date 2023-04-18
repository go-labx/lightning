package main

import (
	"net/http"

	"github.com/go-labx/lightning"
)

func main() {
	app := lightning.NewApp(
		&lightning.Config{
			AppName: "demo",
		},
	)

	app.Get("/ping", func(ctx *lightning.Context) {
		ctx.JSON(http.StatusOK, lightning.Map{
			"message": "pong",
		})
	})

	app.Run()
}
