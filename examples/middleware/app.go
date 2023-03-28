package main

import (
	"fmt"
	"github.com/go-labx/lightning"
)

func main() {
	app := lightning.App()

	// Global scope middleware
	app.Use(func(ctx *lightning.Context, next lightning.Next) {
		fmt.Println("global scope middleware --->")
		next()
		fmt.Println("<--- global scope middleware")
	})

	app.Get("/ping", func(ctx *lightning.Context) {
		ctx.Text("pong")
	})

	// Route scope middleware
	app.Get("/hello", func(ctx *lightning.Context) {
		ctx.Text("hello world")
	}, func(ctx *lightning.Context, next lightning.Next) {
		fmt.Println("Route scope middleware 1 --->")
		next()
		fmt.Println("<--- Route scope middleware 1")
	}, func(ctx *lightning.Context, next lightning.Next) {
		fmt.Println("Route scope middleware 2 --->")
		next()
		fmt.Println("<--- Route scope middleware 2")
	})

	app.Run()
}
