package main

import (
	"fmt"
	"github.com/go-labx/lightning"
	"net/http"
)

func main() {
	app := lightning.NewApp()

	app.NotFoundHandler = func(ctx *lightning.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "/404")
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"code": 404,
			"msg":  fmt.Sprintf("API %s -> %s not found.", ctx.Method, ctx.Path),
		})
	}

	app.Get("/ping", func(ctx *lightning.Context) {
		ctx.JSON(http.StatusOK, map[string]string{
			"message": "pong",
		})
	})

	app.Run()
}
