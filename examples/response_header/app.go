package main

import (
	"github.com/go-labx/lightning"
	"net/http"
)

func main() {
	app := lightning.DefaultApp()

	app.Get("/ping", func(ctx *lightning.Context) {
		// Add multiple headers with the same key "foo"
		ctx.AddHeader("foo", "bar")
		ctx.AddHeader("foo", "baz")
		ctx.AddHeader("foo", "baq")

		// Delete all headers with key "foo"
		ctx.DelHeader("foo")

		// Set a header with key "id" and value "ewh2mime9purchaser4error"
		ctx.SetHeader("id", "ewh2mime9purchaser4error")

		ctx.JSON(http.StatusOK, lightning.Map{
			"message": "pong",
		})
	})

	app.Run()
}