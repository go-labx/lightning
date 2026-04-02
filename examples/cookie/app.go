package main

import (
	"fmt"

	"github.com/go-labx/lightning"
)

func main() {
	app := lightning.DefaultApp()

	app.Get("/ping", func(ctx *lightning.Context) {
		cookie := ctx.Cookie("sid")
		if cookie != "" {
			fmt.Println("sid", cookie)
		}

		cookies := ctx.Cookies()
		for name, value := range cookies {
			fmt.Println(name, value)
		}

		ctx.SetCookie("sid", "sid:xxxxxxxxxx")

		ctx.JSON(lightning.StatusOK, map[string]string{
			"message": "pong",
		})
	})

	app.Run()
}
