package main

import (
	"net/http"

	"github.com/go-labx/lightning"
)

// User struct contains fields for user data
type User struct {
	Name     string `validate:"required" json:"name"`
	Password string `validate:"required,min=8,max=32" json:"password"`
	Email    string `validate:"required,email" json:"email"`
}

func main() {
	app := lightning.DefaultApp()

	app.Post("/validate", func(ctx *lightning.Context) {
		// Create a new User struct
		var user = &User{}

		// Bind and validate the request body to the User struct
		if err := ctx.JSONBodyWithValidate(user); err != nil {
			// If there is an error, return it as JSON
			ctx.JSON(http.StatusOK, lightning.Map{
				"err": err.Error(),
			})
			return
		}

		// If there is no error, return the User struct as JSON
		ctx.JSON(http.StatusOK, lightning.Map{
			"user": user,
		})
	})

	app.Run()
}
