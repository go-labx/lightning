package main

import (
	"fmt"
	"github.com/go-labx/lightning"
)

func main() {
	app := lightning.App()

	app.Use(func(ctx *lightning.Context) {
		fmt.Println("global middleware 1 --->")
		ctx.Next()
		fmt.Println("<--- global middleware 1")
	})
	app.Use(func(ctx *lightning.Context) {
		fmt.Println("global middleware 2 --->")
		ctx.Next()
		fmt.Println("<--- global middleware 2")
	})

	parentGroup := app.Group("/api")

	parentGroup.Use(func(ctx *lightning.Context) {
		fmt.Println("parent group middleware 1 --->")
		ctx.Next()
		fmt.Println("<--- parent group middleware 1")
	})
	parentGroup.Use(func(ctx *lightning.Context) {
		fmt.Println("parent group middleware 2 --->")
		ctx.Next()
		fmt.Println("<--- parent group middleware 2")
	})

	subGroup := parentGroup.Group("/user")
	subGroup.Use(func(ctx *lightning.Context) {
		fmt.Println("sub group middleware 1 --->")
		ctx.Next()
		fmt.Println("<--- sub group middleware 1")
	})

	subGroup.Get("/info", func(ctx *lightning.Context) {
		ctx.JSON(map[string]interface{}{
			"username": "zhangsan",
			"age":      20,
		})
	})

	app.Get("/ping", func(ctx *lightning.Context) {
		ctx.Text("pong")
	})

	app.Run()
}
