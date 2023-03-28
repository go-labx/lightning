package main

import "github.com/go-labx/lightning"

func main() {
	app := lightning.App()

	app.Use(func(ctx *lightning.Context, next lightning.Next) {
		ctx.SetData("session", map[string]interface{}{
			"userId":   123,
			"username": "Jack",
		})
		next()
	})

	app.Get("/", func(ctx *lightning.Context) {
		session := ctx.GetData("session")

		ctx.JSON(session)
	})

	app.Run()
}
