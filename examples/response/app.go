package main

import "github.com/go-labx/lightning"

func main() {
	app := lightning.NewApp()

	app.Get("/text", func(ctx *lightning.Context) {
		ctx.Text("hello world")
	})

	app.Get("/json", func(ctx *lightning.Context) {
		ctx.JSON(map[string]string{
			"message": "pong",
		})
	})

	app.Get("/xml", func(ctx *lightning.Context) {
		type Person struct {
			name string
			age  int
		}
		ctx.XML(Person{
			name: "zhangsan",
			age:  20,
		})
	})

	app.Get("/file", func(ctx *lightning.Context) {
		ctx.File("./LICENSE")
	})

	app.Run()
}
