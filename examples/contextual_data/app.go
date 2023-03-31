package main

import (
	"github.com/go-labx/lightning"
	"net/http"
)

func main() {
	app := lightning.NewApp()

	app.Use(func(ctx *lightning.Context) {
		ctx.SetData("session", map[string]interface{}{
			"userId":   123,
			"username": "Jack",
		})
		ctx.Next()
	})

	app.Get("/", func(ctx *lightning.Context) {
		session := ctx.GetData("session")

		ctx.JSON(http.StatusOK, session)
	})

	app.Run()
}
