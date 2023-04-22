package main

import (
	"net/http"

	"github.com/go-labx/lightning"
)

func main() {
	app := lightning.DefaultApp()

	app.Get("/ping", func(ctx *lightning.Context) {
		// Add multiple headers with the same key "foo"
		ctx.AddHeader("foo", "bar")
		ctx.AddHeader("foo", "baz")
		ctx.AddHeader("foo", "baq")

		// set a header with key "id" and value "ewh2mime9purchaser4error"
		ctx.SetHeader("id", "ewh2mime9purchaser4error")

		// delete all headers with key "bar"
		ctx.DelHeader("bar")

		ctx.JSON(http.StatusOK, lightning.Map{
			"message": "pong",
		})
	})

	app.Run()
}
