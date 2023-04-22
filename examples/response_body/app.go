package main

import (
	"net/http"

	"github.com/go-labx/lightning"
)

type Person struct {
	Name string `xml:"name" json:"name"` // Person's name
	Age  int    `xml:"age" json:"age"`   // Person's age
	City string `xml:"city" json:"city"` // Person's city
}

func main() {
	app := lightning.DefaultApp()

	app.Get("/text", func(ctx *lightning.Context) {
		// Return "hello world" as plain text
		ctx.Text(http.StatusOK, "hello world")
	})

	app.Get("/json", func(ctx *lightning.Context) {
		// Return a Person object as JSON with name "zhangsan", age 20, and city "Hangzhou"
		ctx.JSON(http.StatusOK, &Person{
			Name: "zhangsan",
			Age:  20,
			City: "Hangzhou",
		})
	})

	app.Get("/xml", func(ctx *lightning.Context) {
		// Return a Person object as XML with name "zhangsan", age 20, and city "Hangzhou"
		ctx.XML(http.StatusOK, &Person{
			Name: "zhangsan",
			Age:  20,
			City: "Hangzhou",
		})
	})

	app.Get("/file", func(ctx *lightning.Context) {
		// Serve the README.md file
		err := ctx.File("./README.md")
		if err != nil {
			ctx.Fail(-1, err.Error())
		}
	})

	app.Get("/success", func(ctx *lightning.Context) {
		// Return a successful response with a Person object as JSON with name "zhangsan", age 20, and city "Hangzhou"
		ctx.Success(&Person{
			Name: "zhangsan",
			Age:  20,
			City: "Hangzhou",
		})
	})

	app.Get("/fail", func(ctx *lightning.Context) {
		// Return a failure response with code 9999 and message "network error"
		ctx.Fail(9999, "network error")
	})

	// Start the server on port 6789
	app.Run(":6789")
}
