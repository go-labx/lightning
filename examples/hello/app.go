package main

import "github.com/go-labx/lightning"

func main() {
	app := lightning.App()

	app.Get("/", func(ctx *lightning.Context) {
		ctx.JSON(map[string]string{
			"message": "hello world~",
		})
	})

	app.Run()
}
