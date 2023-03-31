package main

import (
	"github.com/go-labx/lightning"
	"net/http"
)

func main() {
	app := lightning.NewApp()

	app.Get("/text", func(ctx *lightning.Context) {
		ctx.Text(http.StatusOK, "hello world")
	})

	app.Get("/json", func(ctx *lightning.Context) {
		ctx.JSON(http.StatusOK, map[string]string{
			"message": "pong",
		})
	})

	app.Get("/xml", func(ctx *lightning.Context) {
		type Person struct {
			name string
			age  int
		}
		ctx.XML(http.StatusOK, Person{
			name: "zhangsan",
			age:  20,
		})
	})

	app.Get("/file", func(ctx *lightning.Context) {
		ctx.File("./LICENSE")
	})

	app.Run()
}
