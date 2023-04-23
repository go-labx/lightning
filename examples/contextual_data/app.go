package main

import (
	"fmt"

	"github.com/go-labx/lightning"
)

func main() {
	app := lightning.NewApp()

	// Middleware to set session data
	app.Use(func(ctx *lightning.Context) {
		ctx.SetData("session", map[string]interface{}{
			"userId":   123,
			"username": "Jack",
		})
		ctx.Next()
	})

	// Middleware to get session data
	app.Use(func(ctx *lightning.Context) {
		session := ctx.GetData("session")
		// write your logic here...
		fmt.Println(session)

		ctx.Next()
	})

	// Route to handle GET request to /api/user
	app.Get("/api/user", func(ctx *lightning.Context) {
		session := ctx.GetData("session")
		// write your logic here...
		fmt.Println(session)

		ctx.Text(200, "hello world")
	})

	// Run the app
	app.Run()
}
