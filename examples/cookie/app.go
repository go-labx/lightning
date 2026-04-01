package main

import (
	"fmt"

	"github.com/go-labx/lightning"
)

func main() {
	app := lightning.DefaultApp()

	app.Get("/ping", func(ctx *lightning.Context) {
		cookie := ctx.Cookie("sid")
		if cookie != nil {
			fmt.Println(string(cookie.Key()), string(cookie.Value()))
		}

		cookies := ctx.Cookies()
		for _, c := range cookies {
			fmt.Println(string(c.Key()), string(c.Value()))
		}

		ctx.SetCookie("sid", "sid:xxxxxxxxxx")

		ctx.JSON(lightning.StatusOK, map[string]string{
			"message": "pong",
		})
	})

	app.Run()
}
