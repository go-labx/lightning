package main

import (
	"fmt"
	"net/http"

	"github.com/go-labx/lightning"
)

func main() {
	app := lightning.DefaultApp()

	app.Get("/ping", func(ctx *lightning.Context) {
		// get the value of the "sid" cookie
		cookie := ctx.Cookie("sid")
		fmt.Println(cookie)

		// get all cookies
		cookies := ctx.Cookies()
		fmt.Println(cookies)

		// set a new cookie
		ctx.SetCookie("sid", "sid:xxxxxxxxxx")

		// set a custom cookie
		ctx.SetCustomCookie(&http.Cookie{
			Name:  "sessionId",
			Value: "sessionId:xxxxxxxxxx",
			Path:  "/",
		})

		ctx.JSON(http.StatusOK, map[string]string{
			"message": "pong",
		})
	})

	app.Run()
}
