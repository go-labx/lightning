package main

import "github.com/go-labx/lightning"

func main() {
	app := lightning.App()

	app.Get("/ping", func(ctx *lightning.Context) {
		ctx.JSON(map[string]string{
			"message": "pong",
		})
	})

	app.Run()
}
