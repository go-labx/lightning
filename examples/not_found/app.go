package main

import (
	"github.com/go-labx/lightning"
)

func main() {
	app := lightning.NewApp(&lightning.Config{
		NotFoundHandler: func(ctx *lightning.Context) {
			ctx.Text(lightning.StatusNotFound, "custom not found")
		},
	})

	app.Get("/ping", func(ctx *lightning.Context) {
		ctx.JSON(lightning.StatusOK, map[string]string{
			"message": "pong",
		})
	})

	app.Run()
}
