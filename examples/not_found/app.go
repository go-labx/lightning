package main

import (
	"fmt"
	"github.com/go-labx/lightning"
)

func main() {
	app := lightning.NewApp()

	app.NotFoundHandlerFunc = func(ctx *lightning.Context) {
		ctx.JSON(map[string]interface{}{
			"code": 404,
			"msg":  fmt.Sprintf("API %s -> %s not found.", ctx.Method, ctx.Path),
		})
	}

	app.Get("/ping", func(ctx *lightning.Context) {
		ctx.JSON(map[string]string{
			"message": "pong",
		})
	})

	app.Run()
}
