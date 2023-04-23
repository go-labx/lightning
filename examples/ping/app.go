package main

import (
	"net/http"

	"github.com/go-labx/lightning"
)

func main() {
	// Create a new Lightning app
	app := lightning.DefaultApp()

	// Define a GET route for "/ping"
	app.Get("/ping", func(ctx *lightning.Context) {
		// Respond with a JSON message
		ctx.JSON(http.StatusOK, map[string]string{
			"message": "pong",
		})
	})

	// Run the app
	app.Run()
}
