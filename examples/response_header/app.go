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

		// delete all headers with key "foo"
		ctx.DelHeader("foo")

		// set a header with key "id" and value "ewh2mime9purchaser4error"
		ctx.SetHeader("id", "ewh2mime9purchaser4error")

		ctx.JSON(http.StatusOK, lightning.Map{
			"message": "pong",
		})
	})

	app.Run()
}
