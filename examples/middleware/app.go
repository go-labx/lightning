package main

import (
	"fmt"
	"github.com/go-labx/lightning"
	"net/http"
)

func main() {
	// Creating a new instance of lightning application
	app := lightning.NewApp()

	// Adding global middleware to the application
	app.Use(func(ctx *lightning.Context) {
		fmt.Println("global scope middleware 1 --->")
		ctx.Next()
		fmt.Println("<--- global scope middleware 1")
	})
	app.Use(func(ctx *lightning.Context) {
		fmt.Println("global scope middleware 2 --->")
		ctx.Next()
		fmt.Println("<--- global scope middleware 2")
	})
	app.Use(func(ctx *lightning.Context) {
		fmt.Println("global scope middleware 3 --->")
		ctx.Next()
		fmt.Println("<--- global scope middleware 3")
	})

	// Defining a GET route for the root path of the application
	app.Get("/", func(ctx *lightning.Context) {
		ctx.JSON(http.StatusOK, map[string]string{
			"message": "hello world",
		})
	})

	// Defining a GET route for the path "/ping" with a route scoped middleware
	app.Get("/ping", func(ctx *lightning.Context) {
		fmt.Println("route scope middleware --->")
		ctx.Next()
		fmt.Println("<--- route scope middleware")
	}, func(ctx *lightning.Context) {
		ctx.Text(http.StatusOK, "pong")
	})

	app.Run()
}
