package main

import (
	"fmt"
	"github.com/go-labx/lightning"
)

func main() {
	app := lightning.App()

	app.Post("/ping", func(ctx *lightning.Context) {
		rawBody := ctx.RawBody()
		stringBody := ctx.StringBody()

		type Person struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
			City string `json:"city"`
		}
		p := &Person{}
		err := ctx.JSONBody(p)
		if err != nil {
			fmt.Println(err)
			return
		}

		ctx.JSON(map[string]interface{}{
			"rawBody":    rawBody,
			"stringBody": stringBody,
			"jsonBody":   p,
		})
	})

	app.Run()
}
