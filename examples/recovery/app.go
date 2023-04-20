package main

import (
	"github.com/go-labx/lightning"
)

func main() {
	app := lightning.NewApp()

	app.Use(lightning.Logger())
	app.Use(lightning.Recovery(func(ctx *lightning.Context) {
		ctx.Fail(9999, "网络异常，请稍后重试")
	}))

	app.Get("/ping", func(ctx *lightning.Context) {
		panic("panic error")
	})

	app.Run()
}
