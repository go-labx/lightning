package main

import (
	"fmt"
	"github.com/go-labx/lightning"
)

func main() {
	app := lightning.NewApp()

	app.Use(func(ctx *lightning.Context) {
		ctx.SetData("session", map[string]interface{}{
			"userId":   123,
			"username": "Jack",
		})
		ctx.Next()
	})

	app.Use(func(ctx *lightning.Context) {
		session := ctx.GetData("session")
		// write your logic here...
		fmt.Println(session)

		ctx.Next()
	})

	app.Get("/api/user", func(ctx *lightning.Context) {
		session := ctx.GetData("session")
		// write your logic here...
		fmt.Println(session)

		ctx.Text(200, "hello world")
	})

	app.Run()
}
