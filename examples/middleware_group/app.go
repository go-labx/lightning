package main

import (
	"fmt"
	"github.com/go-labx/lightning"
	"net/http"
)

func main() {
	// create a new lightning application
	app := lightning.DefaultApp()

	// create a new group of routes for "/api"
	group := app.Group("/api")

	group.Use(func(ctx *lightning.Context) {
		fmt.Println("routing group middleware 1 --->")
		ctx.Next()
		fmt.Println("<--- routing group middleware 1")
	})
	group.Use(func(ctx *lightning.Context) {
		fmt.Println("routing group middleware 2 --->")
		ctx.Next()
		fmt.Println("<--- routing group middleware 2")
	})

	group.Get("/ping", func(ctx *lightning.Context) {
		ctx.Text(http.StatusOK, "pong")
	})

	app.Run()
}
