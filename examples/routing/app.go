package main

import "github.com/go-labx/lightning"

func main() {
	app := lightning.App()

	// Basics
	app.Post("/api/article", func(ctx *lightning.Context) {
		// implementation
	})

	// Route parameters
	app.Get("/api/article/:articleId", func(ctx *lightning.Context) {
		// implementation
	})
	app.Put("/api/article/:articleId", func(ctx *lightning.Context) {
		// implementation
	})
	app.Patch("/api/article/:articleId/name", func(ctx *lightning.Context) {
		// implementation
	})
	app.Delete("/api/article/:articleId", func(ctx *lightning.Context) {
		// implementation
	})

	// Wildcards
	app.Get("/api/*", func(ctx *lightning.Context) {
		// implementation
	})

	app.Run()
}
