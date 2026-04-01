package main

import (
	"fmt"

	"github.com/go-labx/lightning"
)

func main() {
	app := lightning.NewApp()

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

	app.Get("/", func(ctx *lightning.Context) {
		ctx.JSON(lightning.StatusOK, map[string]string{
			"message": "hello world",
		})
	})

	app.Get("/ping", func(ctx *lightning.Context) {
		fmt.Println("route scope middleware --->")
		ctx.Next()
		fmt.Println("<--- route scope middleware")
	}, func(ctx *lightning.Context) {
		ctx.Text(lightning.StatusOK, "pong")
	})

	app.Run()
}
