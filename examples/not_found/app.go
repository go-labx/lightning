package main

import (
	"net/http"

	"github.com/go-labx/lightning"
)

func main() {
	app := lightning.NewApp(&lightning.Config{
		NotFoundHandler: func(ctx *lightning.Context) {
			ctx.Text(404, "custom not found")
		},
	})

	app.Get("/ping", func(ctx *lightning.Context) {
		ctx.JSON(http.StatusOK, map[string]string{
			"message": "pong",
		})
	})

	app.Run()
}
