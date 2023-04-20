package main

import (
	"fmt"
	"net/http"

	"github.com/go-labx/lightning"
)

func main() {
	// create a new lightning application
	app := lightning.NewApp()

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

	// create a new group of routes for "/api"
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

	// create a new subgroup of routes for "/api/user"
	subGroup := parentGroup.Group("/user")
	subGroup.Use(func(ctx *lightning.Context) {
		fmt.Println("sub group middleware 1 --->")
		ctx.Next()
		fmt.Println("<--- sub group middleware 1")
	})

	// define a GET route for "/api/user/info"
	subGroup.Get("/info", func(ctx *lightning.Context) {
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"username": "zhangsan",
			"age":      20,
		})
	})

	app.Get("/ping", func(ctx *lightning.Context) {
		ctx.Text(http.StatusOK, "pong")
	})

	app.Run()
}
