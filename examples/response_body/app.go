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
		ctx.Text(http.StatusOK, "hello world") // Return "hello world" as plain text
	})

	app.Get("/json", func(ctx *lightning.Context) {
		ctx.JSON(http.StatusOK, lightning.Map{
			"message": "pong", // Return a JSON object with a "message" key and "pong" value
		})
	})

	app.Get("/xml", func(ctx *lightning.Context) {
		ctx.XML(http.StatusOK, Person{
			Name: "zhangsan", // Return a Person object as XML with name "zhangsan", age 20, and city "Hangzhou"
			Age:  20,
			City: "Hangzhou",
		})
	})

	app.Get("/file", func(ctx *lightning.Context) {
		err := ctx.File("./README.md")
		if err != nil {
			ctx.Fail(-1, err.Error())
		}
	})

	app.Get("/success", func(ctx *lightning.Context) {
		ctx.Success(&Person{
			Name: "zhangsan", // Return a successful response with a Person object as JSON with name "zhangsan", age 20, and city "Hangzhou"
			Age:  20,
			City: "Hangzhou",
		})
	})

	app.Get("/fail", func(ctx *lightning.Context) {
		ctx.Fail(9999, "network error") // Return a failure response with code 9999 and message "network error"
	})

	// Start the server on port 6789
	app.Run(":6789")
}
