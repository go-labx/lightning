package main

import (
	"fmt"
	"github.com/go-labx/lightning"
	"html/template"
	"time"
)

func formatDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
	// Create a new Lightning app
	app := lightning.DefaultApp()

	app.Static("./public", "/static")
	app.SetFuncMap(template.FuncMap{
		"formatDate": formatDate,
	})
	app.LoadHTMLGlob("templates/*")

	app.Get("/", func(ctx *lightning.Context) {
		ctx.HTML(200, "index.tmpl", lightning.Map{
			"title":       "Lightning",
			"description": "Lightning is a lightweight and fast web framework for Go. It is designed to be easy to use and highly performant. ⚡️⚡️⚡️",
			"now":         time.Now(),
		})
	})

	// Run the app
	app.Run()
}
